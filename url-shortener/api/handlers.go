package api

import (
	"errors"
	"net/http"
	"sync"
	"url-shortener/config"

	common "url-shortener/db"

	"github.com/gin-gonic/gin"
)

type EntryResponse struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

// Handler struct holds dependencies for HTTP handlers
type Handler struct {
	database  common.Database
	config    *config.Config
	validator *URLValidator
	mu        sync.RWMutex
}

// NewHandler creates a new handler instance
func NewHandler(database common.Database, config *config.Config) *Handler {
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

	entry, found := h.database.GetEntry(id)
	if !found {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Entry not found"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, entry.URL)
}

// HandleGet retrieves the Entry based on the short ID
func (h *Handler) HandleGet(c *gin.Context) {
	id := c.Param("id")
	if !h.validator.IsValidShortID(id) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid short ID format"})
		return
	}

	entry, found := h.database.GetEntry(id)
	if !found {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Entry not found"})
		return
	}

	c.JSON(http.StatusOK, EntryResponse{Id: entry.ID, Url: entry.URL})
}

// HandleCreateShortURL creates a new short URL entry
func (h *Handler) HandleCreateShortURL(c *gin.Context) {
	var entryRequest EntryRequest
	if err := c.ShouldBindJSON(&entryRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	entry, err := h.CreateEntry(entryRequest.URL)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid URL format"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// HandleHealthCheck provides a health check endpoint
func (h *Handler) HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Service is healthy",
	})
}

func (h *Handler) CreateEntry(url string) (common.Entry, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.validator.IsValidURL(url) {
		return common.Entry{}, errors.New("invalid URL format")
	}

	// Normalize the URL before storing
	normalizedURL := h.validator.NormalizeURL(url)
	return h.database.AddEntry(normalizedURL), nil
}
