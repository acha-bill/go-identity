package store

import (
	"go-identity/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIdentity(t *testing.T) {
	iStore := NewIdentity(nil)
	assert.NotNil(t, iStore)
}

func TestIdentity_Get(t *testing.T) {
	d := map[int]*domain.Identity{
		1: {
			FirstName: "John",
		},
	}
	iStore := NewIdentity(d)
	id, err := iStore.Get(1)
	require.NoError(t, err)
	require.NotNil(t, id)
	assert.Equal(t, "John", id.FirstName)

	id, err = iStore.Get(2)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrIdentityNotFound)
}
