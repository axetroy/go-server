package cryptoengine

import (
	"bytes"
	"errors"
	"fmt"
)

// The verification engine links two peers basically.
// It holds the public key and the remote peer public key and the pre-shared key
type VerificationEngine struct {
	publicKey        [keySize]byte // the peer public key
	signingPublicKey [keySize]byte // the peer public signing key => this is not implemented yet, because go does not support Ed25519 signatures yet
}

// This function instantiate the verification engine by leveraging the context
// Basically if a public key of a peer is available locally then it's locaded here
func NewVerificationEngine(context string) (VerificationEngine, error) {

	engine := VerificationEngine{}

	if context == "" {
		return engine, errors.New("Context cannot be empty when initializing the Verification Engine")
	}

	// try to load the public key and if it succeed, then return both the keys
	publicFile := fmt.Sprintf(publicKeySuffixFormat, sanitizeIdentifier(context))
	// if the key exists
	if keyFileExists(publicFile) {
		// try to read it
		public, err := readKey(publicFile, keysFolderPrefixFormat)
		if err != nil {
			// in case of error return it
			return engine, err
		}

		// if we reached here, it means that both the public key
		// existed and was loaded successfull
		engine.publicKey = public
	}

	return engine, nil

}

// This function instantiate the verification engine by passing it the key (at the moment only the public key)
// go nacl crypto does not support Ed25519 signatures yet
func NewVerificationEngineWithKey(publicKey []byte) (VerificationEngine, error) {

	engine := VerificationEngine{}
	var data32 [keySize]byte

	// check the peerPublicKey is not empty (all zeros)
	if bytes.Compare(publicKey[:], emptyKey) == 0 {
		return engine, errors.New("Public key cannot be empty while creating the verification engine")
	}

	total := copy(data32[:], publicKey)
	if total != keySize {
		return engine, KeySizeError
	}

	engine.publicKey = data32
	return engine, nil

}

func (e VerificationEngine) PublicKey() [keySize]byte {
	return e.publicKey
}
