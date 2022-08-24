package handler

import (
	"fmt"
	"go-identity/pkg"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Identity is the identity handler.
type Identity struct {
	log    *log.Logger
	router *mux.Router
	svc    *pkg.Identity
}

// NewIdentity returns a new identity handler.
func NewIdentity(log *log.Logger, router *mux.Router, svc *pkg.Identity) *Identity {
	return &Identity{
		log,
		router,
		svc,
	}
}

// RegisterRoutes registers identity endpoints on the router.
func (i *Identity) RegisterRoutes() {
	i.router.HandleFunc("/identity/{documentHash}", i.getIdentity).Methods(http.MethodGet)
}

func (i *Identity) getIdentity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	documentHash, err := pkg.DocumentHashFromBase64(params["documentHash"])
	if err != nil {
		i.log.Printf("failed to decode document hash: %w\n", err)
		w.WriteHeader(http.StatusBadRequest)
		errResponse := ErrResponse{Message: fmt.Sprintf("failed to decode document hash: %v", err)}
		io.WriteString(w, errResponse.ToJSON())
		return
	}

	signedIdentity, err := i.svc.GetIdentity(documentHash)
	if err != nil {
		i.log.Printf("failed to get signed identity: %w\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		errResponse := ErrResponse{Message: err.Error()}
		io.WriteString(w, errResponse.ToJSON())
		return
	}

	i.log.Println("successfully signed bundle\n")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, signedIdentity.ToJSON())
}
