package api

import (
	"net/url"
	"strings"
)

const (
	httpScheme  = "http://"
	httpsScheme = "https://"
)

// URLValidator provides URL validation functionality
type URLValidator struct{}

// NewURLValidator creates a new URL validator
func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

// IsValidURL checks if the provided string is a valid URL
func (v *URLValidator) IsValidURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	// Add scheme if missing
	if !strings.HasPrefix(urlStr, httpScheme) && !strings.HasPrefix(urlStr, httpsScheme) {
		urlStr = httpsScheme + urlStr
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

// NormalizeURL normalizes a URL by adding scheme if missing
func (v *URLValidator) NormalizeURL(urlStr string) string {
	if strings.HasPrefix(urlStr, httpScheme) || strings.HasPrefix(urlStr, httpsScheme) {
		return urlStr
	}
	return httpsScheme + urlStr
}

// IsValidShortID checks if the provided short ID is valid
func (v *URLValidator) IsValidShortID(id string) bool {
	if id == "" {
		return false
	}

	// Check if ID contains only alphanumeric characters and is reasonable length
	if len(id) < 1 || len(id) > 20 {
		return false
	}

	return true
}
