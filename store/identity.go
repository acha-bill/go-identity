package store

import "go-identity/domain"

// Identity is an identity store.
type Identity struct {
	d map[int]*domain.Identity
}

// NewIdentity returns a new identity store.
func NewIdentity(d map[int]*domain.Identity) *Identity {
	return &Identity{d}
}

// Get returns an identity.
func (i *Identity) Get() *domain.Identity {
	return i.d[0]
}
