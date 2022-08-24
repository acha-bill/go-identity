package pkg

import (
	"go-identity/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBundle_Bytes(t *testing.T) {
	id := &domain.Identity{
		FirstName: "John",
		LastName:  "Doe",
		DOB:       time.Now(),
	}
	bundle := Bundle{
		DocumentHash: DocumentHash{},
		Identity:     id,
	}
	b, err := bundle.Bytes()
	assert.NoError(t, err)
	assert.NotNil(t, b)
}
