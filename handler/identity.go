package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-identity/domain"
	"go-identity/pkg"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Identity is the identity handler.
type Identity struct {
	log    *log.Logger
	router *mux.Router
	svc    *pkg.Identity
}

// GetIdentityReq is the request body for getIdentity.
type GetIdentityReq struct {
	DocumentHash string `json:"documentHash"`
}

// GetIdentityRes is the response for getIdentity.
type GetIdentityRes struct {
	DocumentHash string           `json:"documentHash"`
	Identity     *domain.Identity `json:"identity"`
	Signature    string           `json:"signature"`
}

func (i GetIdentityRes) ToJSON() string {
	b, _ := json.Marshal(i)
	return string(b)
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
	i.router.HandleFunc("/identity", i.getIdentity).Methods(http.MethodPost)
}

func (i *Identity) getIdentity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req GetIdentityReq
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errResponse := ErrResponse{Message: fmt.Sprintf("failed to read req body: %v", err)}
		io.WriteString(w, errResponse.ToJSON())
		return
	}
	err = json.Unmarshal(d, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errResponse := ErrResponse{Message: fmt.Sprintf("failed parse request: %v", err)}
		io.WriteString(w, errResponse.ToJSON())
		return
	}

	documentHash, err := pkg.DocumentHashFromBase64(req.DocumentHash)
	if err != nil {
		i.log.Printf("failed to decode document hash: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errResponse := ErrResponse{Message: fmt.Sprintf("failed to decode document hash: %v", err)}
		io.WriteString(w, errResponse.ToJSON())
		return
	}

	userID, err := strconv.Atoi(r.Header.Get("userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errResponse := ErrResponse{Message: fmt.Sprintf("invalid userID in header: %s", err.Error())}
		io.WriteString(w, errResponse.ToJSON())
		return
	}

	signedIdentity, err := i.svc.GetIdentity(documentHash, userID)
	if err != nil {
		i.log.Printf("failed to get signed identity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errResponse := ErrResponse{Message: fmt.Sprintf("failed to get signed identity: %s", err.Error())}
		io.WriteString(w, errResponse.ToJSON())
		return
	}

	i.log.Println("successfully signed bundle")
	signature := base64.StdEncoding.EncodeToString(signedIdentity.Signature)
	res := GetIdentityRes{
		DocumentHash: documentHash.ToBase64(),
		Identity:     signedIdentity.Identity,
		Signature:    signature,
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, res.ToJSON())
}
