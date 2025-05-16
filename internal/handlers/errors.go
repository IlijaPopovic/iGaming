package handlers

// ErrorResponse represents a standard error response
type ErrorResponse struct {
    Error string `json:"error" example:"error message"`
}

// ValidationErrorResponse represents validation errors
type ValidationErrorResponse struct {
    Error  string            `json:"error" example:"validation failed"`
    Errors map[string]string `json:"errors,omitempty"`
}

type DetailedErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code,omitempty"`
    Details string `json:"details,omitempty"`
}