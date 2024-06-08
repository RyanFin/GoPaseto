package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"` //uuid.UUID is its own type
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiryAt  time.Time `json:"expiry_at"`
}

// find out why the return type has to be a pointer in this case
// Reason 1: By returning a pointer you only only return the memory address which is light weight and constant in size, unlike the entire struct obj which could potentially be huge
// @ Function creates a new payload object
func NewPayload(username string, duration time.Duration) (*Payload, error) {

	// create a new random uuid as a token
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiryAt:  time.Now().Add(duration),
	}

	return payload, nil
}

// @ Function checks the expiry time of the PASETO token
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiryAt) {
		return errors.New("Token has expired...")
	}
	return nil
}
