package api

import (
	"fmt"
	"net/http"

	"better-uptime/config"
	"better-uptime/internal/api/auth"
	db "better-uptime/internal/db/sqlc"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	store       db.Store
	cfg         *config.Config
	rdb         redis.Client
	router      *chi.Mux
	authHandler *auth.Handler
}

type ServerConfig struct {
	Config *config.Config
	Store  db.Store
	// Uncomment these if you need them later
	Redis redis.Client
	// S3Session       s3.S3Interface
	// Firestore       firestore.FirestoreInterface
	// Email           email.EmailInterface
	// CalendarService *calendar.Service
}

// NewServer creates a new API server instance
func NewServer(store db.Store, cfg *config.Config, rdb redis.Client) *Server {

	// Create the server instance first
	server := &Server{
		store: store,
		cfg:   cfg,
		rdb:   rdb,
	}

	// Initialize the auth handler with only required dependencies
	server.authHandler = auth.NewHandler(cfg, store)

	// You can now mount auth routes here like:
	// r.Post("/login", server.authHandler.Login)

	server.router = server.routes()

	return server
}

// Start launches the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.PORT)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	fmt.Println("Server running on", addr)
	return srv.ListenAndServe()
}
