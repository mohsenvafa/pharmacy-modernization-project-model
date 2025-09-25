package helper

import (
	"encoding/json"
	"net/http"
)

type APIResponse[T any] struct {
	Data  *T        `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
	Meta  any       `json:"meta,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func WriteJSON[T any](w http.ResponseWriter, status int, data T, meta any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := APIResponse[T]{Data: &data}
	if meta != nil {
		resp.Meta = meta
	}

	return json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, status int, apiErr APIError) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := APIResponse[struct{}]{Error: &apiErr}
	return json.NewEncoder(w).Encode(resp)
}
