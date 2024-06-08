package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// Maker serves as a JSON web token generator
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// @ Function instantiates a new PasetoMaker generator object and also validates the newly created object
func NewPaseto(symmetricKey string) (*PasetoMaker, error) {
	// symmetric key length check
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("SymmetricKey too short should be: %v", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// @ Function creates a new PASETO token that is encrypted
// @ Params username and duration of the token before its expiry
// @ Returns paseto token string and an error
func (m *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return m.paseto.Encrypt(m.symmetricKey, payload, nil)
}

// @ Function verifies an existing PASETO token, decrypts it and checks its validity
// @ Params token string
// @ Returns a payload object and an error
func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	// create an empty payload that will be populated and returned
	payload := &Payload{}

	// decrypt the existing paseto token using the symmetric key provided
	err := m.paseto.Decrypt(token, m.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	// check the validity of the data
	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
