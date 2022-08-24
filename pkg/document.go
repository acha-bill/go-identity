package pkg

import (
	"crypto/sha512"
	"encoding/base64"
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
