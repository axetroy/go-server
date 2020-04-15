package twofactor

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/sec51/convert"
	"github.com/sec51/convert/bigendian"
	"github.com/sec51/cryptoengine"
	qr "github.com/sec51/qrcode"
	"hash"
	"io"
	"math"
	"net/url"
	"strconv"
	"time"
)

const (
	backoff_minutes = 5 // this is the time to wait before verifying another token
	max_failures    = 3 // total amount of failures, after that the user needs to wait for the backoff time
	counter_size    = 8 // this is defined in the RFC 4226
	message_type    = 0 // this is the message type for the crypto engine
)

var (
	initializationFailedError = errors.New("Totp has not been initialized correctly")
	LockDownError             = errors.New("The verification is locked down, because of too many trials.")
)

// WARNING: The `Totp` struct should never be instantiated manually!
// Use the `NewTOTP` function
type Totp struct {
	key                       []byte             // this is the secret key
	counter                   [counter_size]byte // this is the counter used to synchronize with the client device
	digits                    int                // total amount of digits of the code displayed on the device
	issuer                    string             // the company which issues the 2FA
	account                   string             // usually the suer email or the account id
	stepSize                  int                // by default 30 seconds
	clientOffset              int                // the amount of steps the client is off
	totalVerificationFailures int                // the total amount of verification failures from the client - by default 10
	lastVerificationTime      time.Time          // the last verification executed
	hashFunction              crypto.Hash        // the hash function used in the HMAC construction (sha1 - sha156 - sha512)
}

// This function is used to synchronize the counter with the client
// Offset can be a negative number as well
// Usually it's either -1, 0 or 1
// This is used internally
func (otp *Totp) synchronizeCounter(offset int) {
	otp.clientOffset = offset
}

// Label returns the combination of issuer:account string
func (otp *Totp) label() string {
	return fmt.Sprintf("%s:%s", url.QueryEscape(otp.issuer), otp.account)
}

// Counter returns the TOTP's 8-byte counter as unsigned 64-bit integer.
func (otp *Totp) getIntCounter() uint64 {
	return bigendian.FromUint64(otp.counter)
}

// This function creates a new TOTP object
// This is the function which is needed to start the whole process
// account: usually the user email
// issuer: the name of the company/service
// hash: is the crypto function used: crypto.SHA1, crypto.SHA256, crypto.SHA512
// digits: is the token amount of digits (6 or 7 or 8)
// steps: the amount of second the token is valid
// it autmatically generates a secret key using the golang crypto rand package. If there is not enough entropy the function returns an error
// The key is not encrypted in this package. It's a secret key. Therefore if you transfer the key bytes in the network,
// please take care of protecting the key or in fact all the bytes.
func NewTOTP(account, issuer string, hash crypto.Hash, digits int) (*Totp, error) {

	keySize := hash.Size()
	key := make([]byte, keySize)
	total, err := rand.Read(key)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("TOTP failed to create because there is not enough entropy, we got only %d random bytes", total))
	}

	// sanitize the digits range otherwise it may create invalid tokens !
	if digits < 6 || digits > 8 {
		digits = 8
	}

	return makeTOTP(key, account, issuer, hash, digits)

}

// Private function which initialize the TOTP so that it's easier to unit test it
// Used internnaly
func makeTOTP(key []byte, account, issuer string, hash crypto.Hash, digits int) (*Totp, error) {
	otp := new(Totp)
	otp.key = key
	otp.account = account
	otp.issuer = issuer
	otp.digits = digits
	otp.stepSize = 30 // we set it to 30 seconds which is the recommended value from the RFC
	otp.clientOffset = 0
	otp.hashFunction = hash
	return otp, nil
}

// This function validates the user privided token
// It calculates 3 different tokens. The current one, one before now and one after now.
// The difference is driven by the TOTP step size
// Based on which of the 3 steps it succeeds to validates, the client offset is updated.
// It also updates the total amount of verification failures and the last time a verification happened in UTC time
// Returns an error in case of verification failure, with the reason
// There is a very basic method which protects from timing attacks, although if the step time used is low it should not be necessary
// An attacker can still learn the synchronization offset. This is however irrelevant because the attacker has then 30 seconds to
// guess the code and after 3 failures the function returns an error for the following 5 minutes
func (otp *Totp) Validate(userCode string) error {

	// check Totp initialization
	if err := totpHasBeenInitialized(otp); err != nil {
		return err
	}

	// verify that the token is valid
	if userCode == "" {
		return errors.New("User provided token is empty")
	}

	// check against the total amount of failures
	if otp.totalVerificationFailures >= max_failures && !validBackoffTime(otp.lastVerificationTime) {
		return LockDownError
	}

	if otp.totalVerificationFailures >= max_failures && validBackoffTime(otp.lastVerificationTime) {
		// reset the total verification failures counter
		otp.totalVerificationFailures = 0
	}

	// calculate the sha256 of the user code
	userTokenHash := sha256.Sum256([]byte(userCode))
	userToken := hex.EncodeToString(userTokenHash[:])

	// 1 calculate the 3 tokens
	tokens := make([]string, 3)
	token0Hash := sha256.Sum256([]byte(calculateTOTP(otp, -1)))
	token1Hash := sha256.Sum256([]byte(calculateTOTP(otp, 0)))
	token2Hash := sha256.Sum256([]byte(calculateTOTP(otp, 1)))

	tokens[0] = hex.EncodeToString(token0Hash[:]) // 30 seconds ago token
	tokens[1] = hex.EncodeToString(token1Hash[:]) // current token
	tokens[2] = hex.EncodeToString(token2Hash[:]) // next 30 seconds token

	// if the current time token is valid then, no need to re-sync and return nil
	if tokens[1] == userToken {
		return nil
	}

	// if the 30 seconds ago token is valid then return nil, but re-synchronize
	if tokens[0] == userToken {
		otp.synchronizeCounter(-1)
		return nil
	}

	// if the let's say 30 seconds ago token is valid then return nil, but re-synchronize
	if tokens[2] == userToken {
		otp.synchronizeCounter(1)
		return nil
	}

	otp.totalVerificationFailures++
	otp.lastVerificationTime = time.Now().UTC() // important to have it in UTC

	// if we got here everything is good
	return errors.New("Tokens mismatch.")
}

// Checks the time difference between the function call time and the parameter
// if the difference of time is greater than BACKOFF_MINUTES  it returns true, otherwise false
func validBackoffTime(lastVerification time.Time) bool {
	diff := lastVerification.UTC().Add(backoff_minutes * time.Minute)
	return time.Now().UTC().After(diff)
}

// Basically, we define TOTP as TOTP = HOTP(K, T), where T is an integer
// and represents the number of time steps between the initial counter
// time T0 and the current Unix time.
// T = (Current Unix time - T0) / X, where the
// default floor function is used in the computation.
// For example, with T0 = 0 and Time Step X = 30, T = 1 if the current
// Unix time is 59 seconds, and T = 2 if the current Unix time is
// 60 seconds.
func (otp *Totp) incrementCounter(index int) {
	// Unix returns t as a Unix time, the number of seconds elapsed since January 1, 1970 UTC.
	counterOffset := time.Duration(index*otp.stepSize) * time.Second
	now := time.Now().UTC().Add(counterOffset).Unix()
	otp.counter = bigendian.ToUint64(increment(now, otp.stepSize))
}

// Function which calculates the value of T (see rfc6238)
func increment(ts int64, stepSize int) uint64 {
	T := float64(ts / int64(stepSize)) // TODO: improve this conversions
	n := convert.Round(T)              // round T
	return n                           // convert n to little endian byte array
}

// Generates a new one time password with hmac-(HASH-FUNCTION)
func (otp *Totp) OTP() (string, error) {

	// verify the proper initialization
	if err := totpHasBeenInitialized(otp); err != nil {
		return "", err
	}

	// it uses the index 0, meaning that it calculates the current one
	return calculateTOTP(otp, 0), nil
}

// Private function which calculates the OTP token based on the index offset
// example: 1 * steps or -1 * steps
func calculateTOTP(otp *Totp, index int) string {
	var h hash.Hash

	switch otp.hashFunction {
	case crypto.SHA256:
		h = hmac.New(sha256.New, otp.key)
		break
	case crypto.SHA512:
		h = hmac.New(sha512.New, otp.key)
		break
	default:
		h = hmac.New(sha1.New, otp.key)
		break

	}

	// set the counter to the current step based ont the current time
	// this is necessary to generate the proper OTP
	otp.incrementCounter(index)

	return calculateToken(otp.counter[:], otp.digits, h)

}

func truncateHash(hmac_result []byte, size int) int64 {
	offset := hmac_result[size-1] & 0xf
	bin_code := (uint32(hmac_result[offset])&0x7f)<<24 |
		(uint32(hmac_result[offset+1])&0xff)<<16 |
		(uint32(hmac_result[offset+2])&0xff)<<8 |
		(uint32(hmac_result[offset+3]) & 0xff)
	return int64(bin_code)
}

// this is the function which calculates the HTOP code
func calculateToken(counter []byte, digits int, h hash.Hash) string {

	h.Write(counter)
	hashResult := h.Sum(nil)
	result := truncateHash(hashResult, h.Size())

	mod := int32(result % int64(math.Pow10(digits)))

	fmtStr := fmt.Sprintf("%%0%dd", digits)

	return fmt.Sprintf(fmtStr, mod)
}

// URL returns a suitable URL, such as for the Google Authenticator app
// example: otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example
func (otp *Totp) url() (string, error) {

	// verify the proper initialization
	if err := totpHasBeenInitialized(otp); err != nil {
		return "", err
	}

	secret := base32.StdEncoding.EncodeToString(otp.key)
	u := url.URL{}
	v := url.Values{}
	u.Scheme = "otpauth"
	u.Host = "totp"
	u.Path = otp.label()
	v.Add("secret", secret)
	v.Add("counter", fmt.Sprintf("%d", otp.getIntCounter()))
	v.Add("issuer", otp.issuer)
	v.Add("digits", strconv.Itoa(otp.digits))
	v.Add("period", strconv.Itoa(otp.stepSize))
	switch otp.hashFunction {
	case crypto.SHA256:
		v.Add("algorithm", "SHA256")
		break
	case crypto.SHA512:
		v.Add("algorithm", "SHA512")
		break
	default:
		v.Add("algorithm", "SHA1")
		break
	}
	u.RawQuery = v.Encode()
	return u.String(), nil
}

// QR generates a byte array containing QR code encoded PNG image, with level Q error correction,
// needed for the client apps to generate tokens
// The QR code should be displayed only the first time the user enabled the Two-Factor authentication.
// The QR code contains the shared KEY between the server application and the client application,
// therefore the QR code should be delivered via secure connection.
func (otp *Totp) QR() ([]byte, error) {

	// get the URL
	u, err := otp.url()

	// check for errors during initialization
	// this is already done on the URL method
	if err != nil {
		return nil, err
	}
	code, err := qr.Encode(u, qr.Q)
	if err != nil {
		return nil, err
	}
	return code.PNG(), nil
}

// ToBytes serialises a TOTP object in a byte array
// Sizes:         4        4      N     8       4        4        N         4          N      4     4          4               8                 4
// Format: |total_bytes|key_size|key|counter|digits|issuer_size|issuer|account_size|account|steps|offset|total_failures|verification_time|hashFunction_type|
// hashFunction_type: 0 = SHA1; 1 = SHA256; 2 = SHA512
// The data is encrypted using the cryptoengine library (which is a wrapper around the golang NaCl library)
// TODO:
// 1- improve sizes. For instance the hashFunction_type could be a short.
func (otp *Totp) ToBytes() ([]byte, error) {

	// check Totp initialization
	if err := totpHasBeenInitialized(otp); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	// caluclate the length of the key and create its byte representation
	keySize := len(otp.key)
	keySizeBytes := bigendian.ToInt(keySize) //bigEndianInt(keySize)

	// caluclate the length of the issuer and create its byte representation
	issuerSize := len(otp.issuer)
	issuerSizeBytes := bigendian.ToInt(issuerSize)

	// caluclate the length of the account and create its byte representation
	accountSize := len(otp.account)
	accountSizeBytes := bigendian.ToInt(accountSize)

	totalSize := 4 + 4 + keySize + 8 + 4 + 4 + issuerSize + 4 + accountSize + 4 + 4 + 4 + 8 + 4
	totalSizeBytes := bigendian.ToInt(totalSize)

	// at this point we are ready to write the data to the byte buffer
	// total size
	if _, err := buffer.Write(totalSizeBytes[:]); err != nil {
		return nil, err
	}

	// key
	if _, err := buffer.Write(keySizeBytes[:]); err != nil {
		return nil, err
	}
	if _, err := buffer.Write(otp.key); err != nil {
		return nil, err
	}

	// counter
	counterBytes := bigendian.ToUint64(otp.getIntCounter())
	if _, err := buffer.Write(counterBytes[:]); err != nil {
		return nil, err
	}

	// digits
	digitBytes := bigendian.ToInt(otp.digits)
	if _, err := buffer.Write(digitBytes[:]); err != nil {
		return nil, err
	}

	// issuer
	if _, err := buffer.Write(issuerSizeBytes[:]); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString(otp.issuer); err != nil {
		return nil, err
	}

	// account
	if _, err := buffer.Write(accountSizeBytes[:]); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString(otp.account); err != nil {
		return nil, err
	}

	// steps
	stepsBytes := bigendian.ToInt(otp.stepSize)
	if _, err := buffer.Write(stepsBytes[:]); err != nil {
		return nil, err
	}

	// offset
	offsetBytes := bigendian.ToInt(otp.clientOffset)
	if _, err := buffer.Write(offsetBytes[:]); err != nil {
		return nil, err
	}

	// total_failures
	totalFailuresBytes := bigendian.ToInt(otp.totalVerificationFailures)
	if _, err := buffer.Write(totalFailuresBytes[:]); err != nil {
		return nil, err
	}

	// last verification time
	verificationTimeBytes := bigendian.ToUint64(uint64(otp.lastVerificationTime.Unix()))
	if _, err := buffer.Write(verificationTimeBytes[:]); err != nil {
		return nil, err
	}

	// has_function_type
	switch otp.hashFunction {
	case crypto.SHA256:
		sha256Bytes := bigendian.ToInt(1)
		if _, err := buffer.Write(sha256Bytes[:]); err != nil {
			return nil, err
		}
		break
	case crypto.SHA512:
		sha512Bytes := bigendian.ToInt(2)
		if _, err := buffer.Write(sha512Bytes[:]); err != nil {
			return nil, err
		}
		break
	default:
		sha1Bytes := bigendian.ToInt(0)
		if _, err := buffer.Write(sha1Bytes[:]); err != nil {
			return nil, err
		}
	}

	// encrypt the TOTP bytes
	engine, err := cryptoengine.InitCryptoEngine(otp.issuer)
	if err != nil {
		return nil, err
	}

	// init the message to be encrypted
	message, err := cryptoengine.NewMessage(buffer.String(), message_type)
	if err != nil {
		return nil, err
	}

	// encrypt it
	encryptedMessage, err := engine.NewEncryptedMessage(message)
	if err != nil {
		return nil, err
	}

	return encryptedMessage.ToBytes()

}

// TOTPFromBytes converts a byte array to a totp object
// it stores the state of the TOTP object, like the key, the current counter, the client offset,
// the total amount of verification failures and the last time a verification happened
func TOTPFromBytes(encryptedMessage []byte, issuer string) (*Totp, error) {

	// init the cryptoengine
	engine, err := cryptoengine.InitCryptoEngine(issuer)
	if err != nil {
		return nil, err
	}

	// decrypt the message
	data, err := engine.Decrypt(encryptedMessage)
	if err != nil {
		return nil, err
	}

	// new reader
	reader := bytes.NewReader([]byte(data.Text))

	// otp object
	otp := new(Totp)

	// get the lenght
	lenght := make([]byte, 4)
	_, err = reader.Read(lenght) // read the 4 bytes for the total lenght
	if err != nil && err != io.EOF {
		return otp, err
	}

	totalSize := bigendian.FromInt([4]byte{lenght[0], lenght[1], lenght[2], lenght[3]})
	buffer := make([]byte, totalSize-4)
	_, err = reader.Read(buffer)
	if err != nil && err != io.EOF {
		return otp, err
	}

	// skip the total bytes size
	startOffset := 0
	// read key size
	endOffset := startOffset + 4
	keyBytes := buffer[startOffset:endOffset]
	keySize := bigendian.FromInt([4]byte{keyBytes[0], keyBytes[1], keyBytes[2], keyBytes[3]})

	// read the key
	startOffset = endOffset
	endOffset = startOffset + keySize
	otp.key = buffer[startOffset:endOffset]

	// read the counter
	startOffset = endOffset
	endOffset = startOffset + 8
	b := buffer[startOffset:endOffset]
	otp.counter = [8]byte{b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7]}

	// read the digits
	startOffset = endOffset
	endOffset = startOffset + 4
	b = buffer[startOffset:endOffset]
	otp.digits = bigendian.FromInt([4]byte{b[0], b[1], b[2], b[3]}) //

	// read the issuer size
	startOffset = endOffset
	endOffset = startOffset + 4
	b = buffer[startOffset:endOffset]
	issuerSize := bigendian.FromInt([4]byte{b[0], b[1], b[2], b[3]})

	// read the issuer string
	startOffset = endOffset
	endOffset = startOffset + issuerSize
	otp.issuer = string(buffer[startOffset:endOffset])

	// read the account size
	startOffset = endOffset
	endOffset = startOffset + 4
	b = buffer[startOffset:endOffset]
	accountSize := bigendian.FromInt([4]byte{b[0], b[1], b[2], b[3]})

	// read the account string
	startOffset = endOffset
	endOffset = startOffset + accountSize
	otp.account = string(buffer[startOffset:endOffset])

	// read the steps
	startOffset = endOffset
	endOffset = startOffset + 4
	b = buffer[startOffset:endOffset]
	otp.stepSize = bigendian.FromInt([4]byte{b[0], b[1], b[2], b[3]})

	// read the offset
	startOffset = endOffset
	endOffset = startOffset + 4
	b = buffer[startOffset:endOffset]
	otp.clientOffset = bigendian.FromInt([4]byte{b[0], b[1], b[2], b[3]})

	// read the total failuers
	startOffset = endOffset
	endOffset = startOffset + 4
	b = buffer[startOffset:endOffset]
	otp.totalVerificationFailures = bigendian.FromInt([4]byte{b[0], b[1], b[2], b[3]})

	// read the offset
	startOffset = endOffset
	endOffset = startOffset + 8
	b = buffer[startOffset:endOffset]
	ts := bigendian.FromUint64([8]byte{b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7]})
	otp.lastVerificationTime = time.Unix(int64(ts), 0)

	// read the hash type
	startOffset = endOffset
	endOffset = startOffset + 4
	b = buffer[startOffset:endOffset]
	hashType := bigendian.FromInt([4]byte{b[0], b[1], b[2], b[3]})

	switch hashType {
	case 1:
		otp.hashFunction = crypto.SHA256
		break
	case 2:
		otp.hashFunction = crypto.SHA512
		break
	default:
		otp.hashFunction = crypto.SHA1
	}

	return otp, err
}

// this method checks the proper initialization of the Totp object
func totpHasBeenInitialized(otp *Totp) error {
	if otp == nil || otp.key == nil || len(otp.key) == 0 {
		return initializationFailedError
	}
	return nil
}
