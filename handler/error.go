package handler

import "encoding/json"

type (
	// ErrResponse is an error response..
	ErrResponse struct {
		Message string `json:"message"`
	}
)

// ToJSON returns err response as a json string.
func (e ErrResponse) ToJSON() string {
	d, _ := json.Marshal(e)
	return string(d)
}
