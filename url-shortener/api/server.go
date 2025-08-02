package api

import (
	"log"
	"url/config"
	common "url/db"
	"url/db/in_mem"
	"url/db/redis"

	"github.com/gin-gonic/gin"
)

// Server represents the URL shortener server
type Server struct {
	database common.Database
	router   *gin.Engine
	config   *config.Config
	handler  *Handler
	Port     string
}

// NewServer creates a new server instance
func NewServer() *Server {
	config := config.LoadConfig()

	var database common.Database
	switch config.DatabaseMode {
	case "in_mem":
		log.Println("Using in_mem database")
		database = in_mem.NewInMemoryDatabase()
	case "redis":
		log.Println("Using redis database")
		database = redis.NewRedisDatabase(config)
	}

	router := gin.Default()

	err := router.SetTrustedProxies(config.TrustedProxies)
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	handler := NewHandler(database, config)
	server := &Server{
		database: database,
		router:   router,
		config:   config,
		handler:  handler,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures all the routes for the server
func (s *Server) setupRoutes() {
	// Serve static files
	s.router.Static("/static", s.config.StaticFilesPath)

	// API routes
	s.router.GET("/", s.handler.HandleIndex)
	s.router.GET("/:id", s.handler.HandleRedirect)
	s.router.POST("/", s.handler.HandleCreateShortURL)

	// Health check endpoint
	s.router.GET("/health", s.handler.HandleHealthCheck)
}

// Run starts the server on the specified port
func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}

// GetRouter returns the underlying gin router (useful for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
