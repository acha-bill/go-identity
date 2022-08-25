package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"go-identity/domain"
	"go-identity/pkg"
	"go-identity/store"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIdentity(t *testing.T) {
	h := NewIdentity(nil, nil, nil)
	assert.NotNil(t, h)
}
func TestIdentity_RegisterRoutes(t *testing.T) {
	router := mux.NewRouter()
	h := NewIdentity(nil, router, nil)
	h.RegisterRoutes()

	shouldBeRegistered := map[string]bool{
		"/identity": false,
	}
	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		template, _ := route.GetPathTemplate()
		shouldBeRegistered[template] = true
		return nil
	})

	ok := true
	for _, registered := range shouldBeRegistered {
		if !registered {
			ok = false
			break
		}
	}
	assert.True(t, ok)
}

func TestIdentity_getIdentity(t *testing.T) {
	// setup
	d := map[int]*domain.Identity{
		1: {FirstName: "John"},
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	iStore := store.NewIdentity(d)
	identitySvc := pkg.NewIdentity(iStore)
	router := mux.NewRouter()
	h := NewIdentity(logger, router, identitySvc)
	router.HandleFunc("/identity", h.getIdentity)

	t.Run("invalid docHash", func(t *testing.T) {
		rr := httptest.NewRecorder()
		body, _ := json.Marshal(&GetIdentityReq{DocumentHash: "invalid"})
		req, err := http.NewRequest("POST", "/identity", bytes.NewBuffer(body))
		require.NoError(t, err)
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("valid doc hash but no private key", func(t *testing.T) {
		// make sure there's no private key
		err := os.Setenv("PRIVATE_KEY", "")
		require.NoError(t, err)

		// generate random docHash
		var docHash pkg.DocumentHash
		b := make([]byte, pkg.DocumentHashLength)
		rand.Read(b)
		copy(docHash[:], b)

		rr := httptest.NewRecorder()
		body, _ := json.Marshal(&GetIdentityReq{DocumentHash: docHash.ToBase64()})
		req, err := http.NewRequest("POST", "/identity", bytes.NewBuffer(body))
		require.NoError(t, err)
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("valid doc hash and valid private key", func(t *testing.T) {
		_, err := pkg.GenerateTestPrivateKey()
		require.NoError(t, err)
		// generate random docHash
		var docHash pkg.DocumentHash
		b := make([]byte, pkg.DocumentHashLength)
		rand.Read(b)
		copy(docHash[:], b)

		rr := httptest.NewRecorder()
		body, _ := json.Marshal(&GetIdentityReq{DocumentHash: docHash.ToBase64()})
		req, err := http.NewRequest("POST", "/identity", bytes.NewBuffer(body))
		require.NoError(t, err)
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
