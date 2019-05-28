package cryptoengine

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
	"log"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	nonceSize           = 24 // this is the nonce size, required by NaCl
	keySize             = 32 // this is the nonce size, required by NaCl
	rotateSaltAfterDays = 7  // this is the amount of days the salt is valid - if it crosses this amount a new salt is generated
	tcpVersion          = 0  // this is the current TCP version
)

var (
	KeySizeError           = errors.New(fmt.Sprintf("The provisioned key size is less than: %d\n", keySize))
	KeyNotValidError       = errors.New("The provisioned public key is not valid")
	SaltGenerationError    = errors.New("Could not generate random salt")
	KeyGenerationError     = errors.New("Could not generate random key")
	MessageDecryptionError = errors.New("Could not verify the message. Message has been tempered with!")
	MessageParsingError    = errors.New("Could not parse the Message from bytes")
	messageEmpty           = errors.New("Can not encrypt an empty message")
	whiteSpaceRegEx        = regexp.MustCompile("\\s")
	emptyKey               = make([]byte, keySize)

	// salt for derivating keys
	saltSuffixFormat = "%s_salt.key" // this is the salt file,for instance: sec51_salt.key

	// secret key for symmetric encryption
	secretSuffixFormat = "%s_secret.key" // this is the secret key crypto file, for instance: sec51_secret.key

	// asymmetric keys
	publicKeySuffixFormat = "%s_public.key"  // this is the public key crypto file,for instance: sec51_public.key
	privateSuffixFormat   = "%s_private.key" // this is the private key crypto file,for instance: sec51_priovate.key

	// nonce secret key
	nonceSuffixFormat = "%s_nonce.key" // this is the secret key crypto file used for generating nonces,for instance: sec51_nonce.key
)

// This is the basic object which needs to be instanciated for encrypting messages
// either via public key cryptography or private key cryptography
// The object has the methods necessary to execute all the needed functions to encrypt and decrypt a message, both with symmetric and asymmetric
// crypto
type CryptoEngine struct {
	context          string                   // this is the context used for the key derivation function and for namespacing the key files
	publicKey        [keySize]byte            // cached asymmetric public key
	privateKey       [keySize]byte            // cached asymmetric private key
	secretKey        [keySize]byte            // secret key used for symmetric encryption
	salt             [keySize]byte            // salt for deriving the random nonces
	nonceKey         [keySize]byte            // this key is used for deriving the random nonces. It's different from the privateKey
	mutex            sync.Mutex               // this mutex is used ti make sure that in case the engine is used by multiple thread the pre-shared key is correctly generated
	preSharedKeysMap map[string][keySize]byte // this map holds the combination hash of peer public key as the map key and the preshared key as value used to encrypt
	counter          uint64                   // this is the counter which is appended to the HKDF at each call
	counterMutex     sync.Mutex               // this is the counter mutex for a safe incrementation (TODO: look into atomic)
}

// This function initialize all the necessary information to carry out a secure communication
// either via public key cryptography or secret key cryptography.
// The peculiarity is that the user of this package needs to take care of only one parameter, the communicationIdentifier.
// It defines a unique set of keys between the application and the communicationIdentifier unique end point.
// IMPORTANT: The parameter communicationIdentifier defines several assumptions the code use:
// - it names the secret key files with the comuncationIdentifier prefix. This means that if you want to have different secret keys
//   with different end points, you can differrentiate the key by having different unique communicationIdentifier.
//   It, also, loads the already created keys back in memory based on the communicationIdentifier
// - it does the same with the asymmetric keys
// The communicationIdentifier parameter is URL unescape, trimmed, set to lower case and all the white spaces are replaced with an underscore.
// The publicKey parameter can be nil. In that case the CryptoEngine assumes it has been instanciated for symmetric crypto usage.
func InitCryptoEngine(communicationIdentifier string) (*CryptoEngine, error) {
	// define an error object
	var err error
	// create a new crypto engine object
	ce := new(CryptoEngine)

	// sanitize the communicationIdentifier
	ce.context = sanitizeIdentifier(communicationIdentifier)

	// load or generate the salt
	salt, err := loadSalt(ce.context)
	if err != nil {
		return nil, err
	}
	ce.salt = salt

	// load or generate the corresponding public/private key pair
	ce.publicKey, ce.privateKey, err = loadKeyPairs(ce.context)
	if err != nil {
		return nil, err
	}

	// load or generate the secret key
	secretKey, err := loadSecretKey(ce.context)
	if err != nil {
		return nil, err
	}
	ce.secretKey = secretKey

	// load the nonce key
	nonceKey, err := loadNonceKey(ce.context)
	if err != nil {
		return nil, err
	}
	ce.nonceKey = nonceKey

	// init the map
	ce.preSharedKeysMap = make(map[string][keySize]byte)

	// finally return the CryptoEngine instance
	return ce, nil

}

// this function reads nonceSize random data
func generateSalt() ([keySize]byte, error) {
	var data32 [keySize]byte
	data := make([]byte, keySize)
	_, err := rand.Read(data)
	if err != nil {
		return data32, err
	}
	total := copy(data32[:], data)
	if total != keySize {
		return data32, SaltGenerationError
	}
	return data32, nil
}

// this function reads keySize random data
func generateSecretKey() ([keySize]byte, error) {
	var data32 [keySize]byte
	data := make([]byte, keySize)
	_, err := rand.Read(data)
	if err != nil {
		return data32, err
	}
	total := copy(data32[:], data[:keySize])
	if total != keySize {
		return data32, KeyGenerationError
	}
	return data32, nil
}

// load the salt random bytes from the id_salt.key
// if the file does not exist, create a new one
// if the file is older than N days (default 2) generate a new one and overwrite the old
// TODO: rotate the salt file
func loadSalt(id string) ([keySize]byte, error) {

	var salt [keySize]byte

	saltFile := fmt.Sprintf(saltSuffixFormat, id)
	if keyFileExists(saltFile) {
		return readKey(saltFile, keysFolderPrefixFormat)
	}

	// generate the random salt
	salt, err := generateSalt()
	if err != nil {
		return salt, err
	}

	// write the salt to the file with its prefix
	if err := writeKey(saltFile, keysFolderPrefixFormat, salt[:]); err != nil {
		return salt, err
	}

	// return the salt and no error
	return salt, nil
}

// load the key random bytes from the id_secret.key
// if the file does not exist, create a new one
func loadSecretKey(id string) ([keySize]byte, error) {

	var key [keySize]byte

	keyFile := fmt.Sprintf(secretSuffixFormat, id)
	if keyFileExists(keyFile) {
		return readKey(keyFile, keysFolderPrefixFormat)
	}

	// generate the random salt
	key, err := generateSecretKey()
	if err != nil {
		return key, err
	}

	// write the salt to the file with its prefix
	if err := writeKey(keyFile, keysFolderPrefixFormat, key[:]); err != nil {
		return key, err
	}

	// return the salt and no error
	return key, nil
}

// load the nonce key random bytes from the id_nonce.key
// if the file does not exist, create a new one
func loadNonceKey(id string) ([keySize]byte, error) {

	var nonceKey [keySize]byte

	nonceFile := fmt.Sprintf(nonceSuffixFormat, id)
	if keyFileExists(nonceFile) {
		return readKey(nonceFile, keysFolderPrefixFormat)
	}

	// generate the random salt
	nonceKey, err := generateSecretKey()
	if err != nil {
		return nonceKey, err
	}

	// write the salt to the file with its prefix
	if err := writeKey(nonceFile, keysFolderPrefixFormat, nonceKey[:]); err != nil {
		return nonceKey, err
	}

	// return the salt and no error
	return nonceKey, nil
}

// load the key pair, public and private keys, the id_public.key, id_private.key
// if the files do not exist, create them
// Returns the publicKey, privateKey, error
func loadKeyPairs(id string) ([keySize]byte, [keySize]byte, error) {

	var private [keySize]byte
	var public [keySize]byte
	var err error

	// try to load the private key
	privateFile := fmt.Sprintf(privateSuffixFormat, id)
	if keyFileExists(privateFile) {
		if private, err = readKey(privateFile, keysFolderPrefixFormat); err != nil {
			return public, private, err
		}
	}
	// try to load the public key and if it succeed, then return both the keys
	publicFile := fmt.Sprintf(publicKeySuffixFormat, id)
	if keyFileExists(publicFile) {
		if public, err = readKey(publicFile, keysFolderPrefixFormat); err != nil {
			return public, private, err
		}

		// if we reached here, it means that both the private and the public key
		// existed and loaded successfully
		return public, private, err
	}

	// if we reached here then, we need to cerate the key pair
	tempPublic, tempPrivate, err := box.GenerateKey(rand.Reader)

	// check for errors first, otherwise continue and store the keys to files
	if err != nil {
		return public, private, err
	}
	// dereference the pointers
	public = *tempPublic
	private = *tempPrivate

	// write the public key first
	if err := writeKey(publicFile, keysFolderPrefixFormat, public[:]); err != nil {
		return public, private, err
	}

	// write the private
	if err := writeKey(privateFile, keysFolderPrefixFormat, private[:]); err != nil {
		// delete the public key, otherwise we remain in an unwanted state
		// the delete can fail as well, therefore we print an error
		if err := deleteFile(publicFile); err != nil {
			log.Printf("[SEVERE] - The private key for asymmetric encryption, %s, failed to be persisted. \nWhile trying to cleanup also the public key previosuly stored, %s, the operation failed as well.\nWe are now in an unrecoverable state.Please delete both files manually: %s - %s", privateFile, publicFile, privateFile, publicFile)
			return public, private, err
		}
		return public, private, err
	}

	// return the data
	return public, private, err

}

// Sanitizes the input of the communicationIdentifier
// The input is URL unescape, trimmed, set to lower case and all the white spaces are replaced with an underscore.
// TODO: evaluate the QueryUnescape error
func sanitizeIdentifier(id string) string {
	// unescape in case it;s URL encoded
	unescaped, _ := url.QueryUnescape(id)
	// trim white spaces
	trimmed := strings.TrimSpace(unescaped)
	// make lower case
	lowered := strings.ToLower(trimmed)
	// replace the white spaces with _
	cleaned := whiteSpaceRegEx.ReplaceAllLiteralString(lowered, "_")
	return cleaned
}

func (engine *CryptoEngine) fetchAndIncrement() string {
	engine.counterMutex.Lock()
	defer engine.counterMutex.Unlock()

	// first read the current value
	// reset the counter
	if engine.counter == math.MaxUint64 {
		engine.counter = 0
	}

	// convert the counter to string
	counterString := strconv.FormatUint(engine.counter, 10)

	// increment the counter
	engine.counter += 1

	return counterString
}

// Gives access to the public key
func (engine *CryptoEngine) PublicKey() []byte {
	return engine.publicKey[:]
}

// This method accepts a message , then encrypts its Version+Type+Text using a symmetric key
func (engine *CryptoEngine) NewEncryptedMessage(msg message) (EncryptedMessage, error) {

	m := EncryptedMessage{}

	// derive nonce
	nonce, err := deriveNonce(engine.nonceKey, engine.salt, engine.context, engine.fetchAndIncrement())
	if err != nil {
		return m, err
	}

	m.nonce = nonce

	encryptedData := secretbox.Seal(nil, msg.toBytes(), &m.nonce, &engine.secretKey)

	// assign the encrypted data to the message
	m.data = encryptedData

	// calculate the overall size of the message
	m.length = uint64(len(m.data) + len(m.nonce) + 8)

	return m, nil

}

// This method accepts the message as byte slice and the public key of the receiver of the messae,
// then encrypts it using the asymmetric key public key.
// If the public key is not privisioned and does not have the required length of 32 bytes it raises an exception.
func (engine *CryptoEngine) NewEncryptedMessageWithPubKey(msg message, verificationEngine VerificationEngine) (EncryptedMessage, error) {

	encryptedMessage := EncryptedMessage{}

	// get the peer public key
	peerPublicKey := verificationEngine.PublicKey()

	// check the size of the peerPublicKey
	if len(peerPublicKey) != keySize {
		return encryptedMessage, KeyNotValidError
	}

	// check the peerPublicKey is not empty (all zeros)
	if bytes.Compare(peerPublicKey[:], emptyKey) == 0 {
		return encryptedMessage, KeyNotValidError
	}

	// derive nonce
	nonce, err := deriveNonce(engine.nonceKey, engine.salt, engine.context, engine.fetchAndIncrement())
	if err != nil {
		return encryptedMessage, err
	}

	// set the nonce to the encrypted message
	encryptedMessage.nonce = nonce

	// calculate the hash of the peer public key
	sha224String := fmt.Sprintf("%x", sha256.Sum224(peerPublicKey[:]))

	// lock the mutex
	engine.mutex.Lock()

	// check if the pre sgared key is already present in the map
	if preSharedKey, ok := engine.preSharedKeysMap[sha224String]; ok { // means the key is there
		// unlock the mutex
		engine.mutex.Unlock()

		// encrypt with the pre-computed key
		encryptedData := box.SealAfterPrecomputation(nil, msg.toBytes(), &nonce, &preSharedKey)

		// assign the encrypted data to the message
		encryptedMessage.data = encryptedData

	} else { // means the key is not there

		// generate the key
		// init the buffer
		preSharedKey = [keySize]byte{}

		// precompute the share key
		box.Precompute(&preSharedKey, &peerPublicKey, &engine.privateKey)

		// assign it to the map
		engine.preSharedKeysMap[sha224String] = preSharedKey

		// unlock the mutex once the map has the sharedKey set
		engine.mutex.Unlock()

		// encrypt with the pre-computed key
		encryptedData := box.SealAfterPrecomputation(nil, msg.toBytes(), &nonce, &preSharedKey)

		// assign the encrypted data to the message
		encryptedMessage.data = encryptedData
	}

	// calculate the size of the message
	encryptedMessage.length = uint64(len(encryptedMessage.data) + len(encryptedMessage.nonce) + 8)

	return encryptedMessage, nil

}

// This method is used to decrypt messages where symmetrci encryption is used
func (engine *CryptoEngine) Decrypt(encryptedBytes []byte) (*message, error) {

	var err error
	msg := new(message)

	// convert the bytes to an encrypted message
	encryptedMessage, err := encryptedMessageFromBytes(encryptedBytes)
	if err != nil {
		return nil, err
	}

	decryptedMessageBytes, valid := secretbox.Open(nil, encryptedMessage.data, &encryptedMessage.nonce, &engine.secretKey)

	// if the verification failed
	if !valid {
		return nil, MessageDecryptionError
	}

	// means we successfully managed to decrypt
	msg, err = messageFromBytes(decryptedMessageBytes)
	return msg, nil

}

// This method is used to decrypt messages where symmetrci encryption is used
func (engine *CryptoEngine) DecryptWithPublicKey(encryptedBytes []byte, verificationEngine VerificationEngine) (*message, error) {

	var err error

	// get the peer public key
	peerPublicKey := verificationEngine.PublicKey()

	// convert the bytes to an encrypted message
	encryptedMessage, err := encryptedMessageFromBytes(encryptedBytes)
	if err != nil {
		return nil, err
	}

	// Make sure the key has a valid size
	if len(peerPublicKey) < keySize {
		return nil, KeyNotValidError
	}

	// calculate the hash of the peer public key
	sha224String := fmt.Sprintf("%x", sha256.Sum224(peerPublicKey[:]))

	// lock the mutex
	engine.mutex.Lock()

	// check if the pre sgared key is already present in the map
	if preSharedKey, ok := engine.preSharedKeysMap[sha224String]; ok { // means the key is there
		// unlock the mutex
		engine.mutex.Unlock()

		messageBytes, err := decryptWithPreShared(preSharedKey, encryptedMessage)
		if err != nil {
			return nil, err
		}
		return messageFromBytes(messageBytes)

	} else {
		// otherwise decrypt with the standard box open function
		messageBytes, valid := box.Open(nil, encryptedMessage.data, &encryptedMessage.nonce, &peerPublicKey, &engine.privateKey)
		if !valid {
			return nil, MessageDecryptionError
		}
		return messageFromBytes(messageBytes)
	}

}

func decryptWithPreShared(preSharedKey [keySize]byte, m EncryptedMessage) ([]byte, error) {
	if decryptedMessage, valid := box.OpenAfterPrecomputation(nil, m.data, &m.nonce, &preSharedKey); !valid {
		return nil, MessageDecryptionError
	} else {
		return decryptedMessage, nil
	}
}
