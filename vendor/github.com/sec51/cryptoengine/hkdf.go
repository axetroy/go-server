package cryptoengine

import (
	"crypto/sha256"
	"errors"
	"golang.org/x/crypto/hkdf"
	"io"
)

// IMPORTANT !!!
// If someone changes the hash function, then the salt needs to have the exactly same lenght!
// So be careful when touching this.
func deriveNonce(masterKey [keySize]byte, salt [keySize]byte, context string, counterValue string) ([nonceSize]byte, error) {
	var data24 [nonceSize]byte
	// Underlying hash function to use
	hash := sha256.New

	// Create the key derivation function
	hkdf := hkdf.New(hash, masterKey[:], salt[:], []byte(context+counterValue))
	// Generate the required keys
	key := make([]byte, nonceSize)
	n, err := io.ReadFull(hkdf, key)
	if n != len(key) || err != nil {
		return data24, err
	}

	total := copy(data24[:], key[:nonceSize])
	if total != nonceSize {
		return data24, errors.New("Could not derive a nonce.")
	}
	return data24, nil

}
