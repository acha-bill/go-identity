package pkg

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"go-identity/domain"
	"go-identity/store"
	"go-identity/utils"
)

type (
	// Identity is the identity service.
	Identity struct {
		s store.Store
	}

	SignedIdentity struct {
		Bundle
		Signature []byte
	}
)

// ToJSON returns the signed identity as a json string.
func (s *SignedIdentity) ToJSON() string {
	d, _ := json.Marshal(s)
	return string(d)
}

// NewIdentity returns a new identity service.
func NewIdentity(s store.Store) *Identity {
	return &Identity{s}
}

// GetIdentity returns an identity from the identity store.
func (i *Identity) GetIdentity(docHash DocumentHash, userID int) (*SignedIdentity, error) {
	id, err := i.s.Get(userID)
	if err != nil {
		return nil, err
	}
	bundle := i.bundle(docHash, id)
	return i.sign(bundle)
}

// bundle returns the document hash and user identity together.
func (i *Identity) bundle(docHash DocumentHash, id *domain.Identity) Bundle {
	return Bundle{
		DocumentHash: docHash,
		Identity:     id,
	}
}

// sign returns the signed bundle.
func (i *Identity) sign(bundle Bundle) (*SignedIdentity, error) {
	msg, err := bundle.Bytes()
	if err != nil {
		return nil, err
	}

	// hash before signing.
	d := sha512.Sum512(msg)

	pemString := utils.GetEnv("PRIVATE_KEY")
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("invalid signature block")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, d[:])
	signedIdentity := SignedIdentity{
		Bundle:    bundle,
		Signature: signature,
	}
	return &signedIdentity, nil
}
