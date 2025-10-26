package helper

import (
	"encoding/json"
	"net/http"

	"pharmacy-modernization-project-model/internal/bind"
)

// Respond400 sends a 400 Bad Request response with field errors
func Respond400(w http.ResponseWriter, ferrs []bind.FieldError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":   "bad_request",
		"message": "invalid request parameters",
		"details": ferrs,
	})
}

// Respond422 sends a 422 Unprocessable Entity response with field errors
func Respond422(w http.ResponseWriter, ferrs []bind.FieldError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":   "validation_failed",
		"message": "request validation failed",
		"details": ferrs,
	})
}

// WriteOK sends a 200 OK response with the provided data
func WriteOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// WriteCreated sends a 201 Created response with the provided data
func WriteCreated(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteNoContent sends a 204 No Content response
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// WriteInternalError sends a 500 Internal Server Error response
func WriteInternalError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusInternalServerError, APIError{
		Code:    "internal_server_error",
		Message: message,
	})
}

// WriteNotFound sends a 404 Not Found response
func WriteNotFound(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusNotFound, APIError{
		Code:    "not_found",
		Message: message,
	})
}

// WriteUnauthorized sends a 401 Unauthorized response
func WriteUnauthorized(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusUnauthorized, APIError{
		Code:    "unauthorized",
		Message: message,
	})
}

// WriteForbidden sends a 403 Forbidden response
func WriteForbidden(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusForbidden, APIError{
		Code:    "forbidden",
		Message: message,
	})
}
