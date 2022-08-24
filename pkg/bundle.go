package pkg

import (
	"bytes"
	"encoding/gob"
	"go-identity/domain"
)

// Bundle is bundle of a document hash and identity.
type Bundle struct {
	DocumentHash DocumentHash
	Identity     *domain.Identity
}

// Bytes encodes the bundle and returns its bytes.
func (b Bundle) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(b); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
