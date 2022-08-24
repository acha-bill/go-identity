package middleware

import (
	"go-identity/handler"
	"io"
	"log"
	"net/http"
)

// Recovery is a recovery middleware.
type Recovery struct {
	log *log.Logger
}

// NewRecovery returns a new recovery middleware.
func NewRecovery(log *log.Logger) *Recovery {
	return &Recovery{log}
}

func (m *Recovery) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				m.log.Println(err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				res := handler.ErrResponse{Message: "an error occurred"}
				io.WriteString(w, res.ToJSON())
			}
		}()
		next.ServeHTTP(w, r)
	})
}
