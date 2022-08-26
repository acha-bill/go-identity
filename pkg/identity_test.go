package pkg

import (
	"crypto/ecdsa"
	"crypto/sha512"
	"encoding/json"
	"go-identity/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockIdentity is a identity store.
type MockIdentityStore struct {
	mock.Mock
}

// Get returns an identity from the store.
func (s *MockIdentityStore) Get(int) (*domain.Identity, error) {
	args := s.Called()
	return args.Get(0).(*domain.Identity), args.Error(1)
}

func TestNewIdentity(t *testing.T) {
	svc := NewIdentity(nil)
	assert.NotNil(t, svc)
}

func TestSignedIdentity_ToJSON(t *testing.T) {
	si := &SignedIdentity{
		Bundle: Bundle{
			DocumentHash: DocumentHash{},
			Identity: &domain.Identity{
				FirstName: "John",
				LastName:  "Doe",
				DOB:       time.Now(),
			},
		},
		Signature: []byte("signature"),
	}
	js := si.ToJSON()
	var sii SignedIdentity
	err := json.Unmarshal([]byte(js), &sii)
	assert.NoError(t, err)
	assert.Equal(t, "John", sii.Identity.FirstName)
}

func TestIdentity_GetIdentity(t *testing.T) {
	s := &MockIdentityStore{}
	svc := NewIdentity(s)
	_, err := GenerateTestPrivateKey()
	require.NoError(t, err)

	i := &domain.Identity{FirstName: "John"}
	s.On("Get").Return(i, nil)
	var docHash DocumentHash
	signedIdentity, err := svc.GetIdentity(docHash, 1)
	require.NoError(t, err)
	assert.Equal(t, "John", signedIdentity.Identity.FirstName)
	assert.Greater(t, len(signedIdentity.Signature), 1)
}

func TestIdentity_bundle(t *testing.T) {
	svc := NewIdentity(nil)
	var documentHash DocumentHash
	i := &domain.Identity{FirstName: "John"}
	bundle := svc.bundle(documentHash, i)
	assert.Equal(t, "John", bundle.Identity.FirstName)
}

func TestIdentity_sign(t *testing.T) {
	svc := NewIdentity(nil)
	bundle := Bundle{
		DocumentHash: DocumentHash{},
		Identity:     &domain.Identity{FirstName: "John"},
	}

	privateKey, err := GenerateTestPrivateKey()
	require.NoError(t, err)

	// sign the bundle
	i, err := svc.sign(bundle)
	require.NoError(t, err)
	require.NotNil(t, i)

	// verify sig
	b, err := bundle.Bytes()
	require.NoError(t, err)
	hash := sha512.Sum512(b)
	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], i.Signature)
	assert.True(t, valid)
}
