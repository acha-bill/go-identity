package store

import "go-identity/domain"

// Store is a store interface.
type Store interface {
	Get() *domain.Identity
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
func (i *Identity) Get() *domain.Identity {
	// dummy: return the first value in the map
	for _, v := range i.d {
		return v
	}
	return nil
}

var _ Store = &Identity{}
