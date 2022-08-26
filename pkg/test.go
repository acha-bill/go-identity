package pkg

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// GenerateTestPrivateKey generates a test private key and adds it to os env.
func GenerateTestPrivateKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	b, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = pem.Encode(buf, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	if err != nil {
		return nil, err
	}
	key := buf.String()

	// set key in env.
	err = os.Setenv("PRIVATE_KEY", key)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
