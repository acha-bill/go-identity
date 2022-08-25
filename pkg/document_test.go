package pkg

import (
	"encoding/base64"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentHash_ToBase64(t *testing.T) {
	var docHash DocumentHash
	b64 := docHash.ToBase64()
	_, err := base64.StdEncoding.DecodeString(b64)
	assert.NoError(t, err)
}

func TestDocumentHashFromBase64(t *testing.T) {
	t.Run("fails on invalid base64", func(t *testing.T) {
		s := "invalid"
		docHash, err := DocumentHashFromBase64(s)
		assert.Error(t, err)
		assert.Equal(t, DocumentHash{}, docHash)
	})

	t.Run("fails on invalid length", func(t *testing.T) {
		s := "c29tZSBkYXRhIHdpdGggACBhbmQg77u/"
		docHash, err := DocumentHashFromBase64(s)
		assert.Error(t, err)
		assert.ErrorIs(t, ErrInvalidDocumentHash, err)
		target := DocumentHash{}
		assert.Equal(t, target, docHash)
	})

	t.Run("valid string", func(t *testing.T) {
		d := make([]byte, DocumentHashLength)
		rand.Read(d)
		s := base64.StdEncoding.EncodeToString(d)
		docHash, err := DocumentHashFromBase64(s)
		assert.NoError(t, err)

		var target DocumentHash
		copy(target[:], d)
		assert.Equal(t, target, docHash)
	})
}
