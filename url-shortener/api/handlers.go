package api

import (
	"net/http"

	common "url/db"

	"github.com/gin-gonic/gin"
)

// Handler struct holds dependencies for HTTP handlers
type Handler struct {
	database  common.Database
	config    *Config
	validator *URLValidator
}

// NewHandler creates a new handler instance
func NewHandler(database common.Database, config *Config) *Handler {
	return &Handler{
		database:  database,
		config:    config,
		validator: NewURLValidator(),
	}
}

// HandleIndex serves the main page
func (h *Handler) HandleIndex(c *gin.Context) {
	c.File(h.config.IndexFilePath)
}

// HandleRedirect redirects to the original URL based on the short ID
func (h *Handler) HandleRedirect(c *gin.Context) {
	id := c.Param("id")
	if !h.validator.IsValidShortID(id) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid short ID format"})
		return
	}

	entry, err := h.database.GetEntry(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Entry not found"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, entry.URL)
}

// HandleCreateShortURL creates a new short URL entry
func (h *Handler) HandleCreateShortURL(c *gin.Context) {
	var entryRequest EntryRequest
	if err := c.ShouldBindJSON(&entryRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	if !h.validator.IsValidURL(entryRequest.URL) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid URL format"})
		return
	}

	// Normalize the URL before storing
	normalizedURL := h.validator.NormalizeURL(entryRequest.URL)
	entry := h.database.AddEntry(normalizedURL)
	c.JSON(http.StatusOK, entry)
}

// HandleHealthCheck provides a health check endpoint
func (h *Handler) HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Service is healthy",
	})
}
