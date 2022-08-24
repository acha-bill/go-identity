package pkg

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"go-identity/domain"
	"go-identity/store"
	"go-identity/utils"
)

type (
	// Identity is the identity service.
	Identity struct {
		s *store.Identity
	}

	SignedIdentity struct {
		Bundle    `json:"bundle"`
		Signature []byte `json:"signature"`
	}
)

// ToJSON returns the signed identity as a json string.
func (s *SignedIdentity) ToJSON() string {
	d, _ := json.Marshal(s)
	return string(d)
}

// NewIdentity returns a new identity service.
func NewIdentity(s *store.Identity) *Identity {
	return &Identity{s}
}

// GetIdentity returns an identity from the identity store.
func (i *Identity) GetIdentity(docHash DocumentHash) (*SignedIdentity, error) {
	id := i.s.Get()
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
	d, err := bundle.Bytes()
	if err != nil {
		return nil, err
	}

	privateKeyStr := utils.GetEnv("PRIVATE_KEY", "secret")
	privateKey, err := x509.ParseECPrivateKey([]byte(privateKeyStr))
	if err != nil {
		return nil, err
	}

	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, d)
	signedIdentity := SignedIdentity{
		Bundle:    bundle,
		Signature: signature,
	}
	return &signedIdentity, nil
}
