package api

// EntryRequest represents the request body for creating a new URL entry
type EntryRequest struct {
	URL string `json:"url" binding:"required"`
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
