package store

import (
	"errors"
	"go-identity/domain"
)

// ErrIdentityNotFound is returned when the user ID is not found in the store.
var ErrIdentityNotFound = errors.New("user identity not found")

// Store is a store interface.
type Store interface {
	Get(ID int) (*domain.Identity, error)
}

// Identity is an identity store.
type Identity struct {
	d map[int]*domain.Identity
}

// NewIdentity returns a new identity store.
func NewIdentity(d map[int]*domain.Identity) *Identity {
	return &Identity{d}
}

// Get returns an identity.
func (i *Identity) Get(ID int) (*domain.Identity, error) {
	id, ok := i.d[ID]
	if !ok {
		return nil, ErrIdentityNotFound
	}
	return id, nil
}

var _ Store = &Identity{}
