package http

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error in API responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// JSON writes a JSON response with the given status code and data
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails after headers are written, we can't send a proper error response.
		// Logging is the best we can do here.
		// In production, consider using a logger if available.
		_ = err
	}
}

// Error writes an error response with the given status code and error
func Error(w http.ResponseWriter, status int, err error) {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: err.Error(),
	}
	JSON(w, status, response)
}

// MethodNotAllowed writes a method not allowed error response
func MethodNotAllowed(w http.ResponseWriter) {
	response := ErrorResponse{
		Error:   "Method not allowed",
		Message: "Method not allowed",
	}
	JSON(w, http.StatusMethodNotAllowed, response)
}

// NotFound writes a not found error response
func NotFound(w http.ResponseWriter) {
	response := ErrorResponse{
		Error:   "Not found",
		Message: "Not found",
	}
	JSON(w, http.StatusNotFound, response)
}

// RequireMethod validates that the request method matches the required method.
// If it doesn't match, it writes a method not allowed response and returns false.
// If it matches, it returns true.
func RequireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		MethodNotAllowed(w)
		return false
	}
	return true
}
