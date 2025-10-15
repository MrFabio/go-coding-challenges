package api

import (
	"log"
	"url-shortener/config"
	common "url-shortener/db"
	"url-shortener/db/in_mem"
	"url-shortener/db/redis"

	"github.com/gin-gonic/gin"
)

// Server represents the URL shortener grpcServer
type Server struct {
	database common.Database
	router   *gin.Engine
	config   *config.Config
	handler  *Handler
	Port     string
}

// NewServer creates a new grpcServer instance
func NewServer() *Server {
	serverConfig := config.LoadConfig()

	var database common.Database
	switch serverConfig.DatabaseMode {
	case "in_mem":
		log.Println("Using in_mem database")
		database = in_mem.NewInMemoryDatabase()
	case "redis":
		log.Println("Using redis database")
		database = redis.NewRedisDatabase(serverConfig)
	}

	router := gin.Default()

	err := router.SetTrustedProxies(serverConfig.TrustedProxies)
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	handler := NewHandler(database, serverConfig)
	server := &Server{
		database: database,
		router:   router,
		config:   serverConfig,
		handler:  handler,
	}

	server.setupRoutes()

	// setup grpc grpcServer
	go registerGrpcServer(handler)

	return server
}

// setupRoutes configures all the routes for the GrpcServer
func (s *Server) setupRoutes() {
	// Serve static files
	s.router.Static("/static", s.config.StaticFilesPath)

	// API routes
	s.router.GET("/", s.handler.HandleIndex)
	s.router.GET("/:id", s.handler.HandleRedirect)
	s.router.GET("/api/:id", s.handler.HandleGet)
	s.router.POST("/", s.handler.HandleCreateShortURL)

	// Health check endpoint
	s.router.GET("/health", s.handler.HandleHealthCheck)
}

// Run starts the grpcServer on the specified port
func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}

// GetRouter returns the underlying gin router (useful for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
