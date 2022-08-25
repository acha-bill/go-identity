package pkg

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
)

const (
	DocumentHashLength = sha512.Size
)

var (
	// ErrInvalidDocumentHash is returned when the base64 hash length is invalid.
	ErrInvalidDocumentHash = errors.New("invalid document hash")
)

// DocumentHash is the hash of a document.
type DocumentHash [DocumentHashLength]byte

// ToBase64 returns the base64 encoding of the document hash.
func (d DocumentHash) ToBase64() string {
	return base64.StdEncoding.EncodeToString(d[:])
}

func (d DocumentHash) MarshalJSON() ([]byte, error) {
	var b []byte
	copy(b, d[:])
	return json.Marshal(b)
}

// DocumentHashFromBase64 converts a base64 string to a valid document hash.
func DocumentHashFromBase64(s string) (DocumentHash, error) {
	var documentHash DocumentHash
	d, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return documentHash, err
	}
	if len(d) != DocumentHashLength {
		return documentHash, ErrInvalidDocumentHash
	}
	copy(documentHash[:], d)
	return documentHash, nil
}
