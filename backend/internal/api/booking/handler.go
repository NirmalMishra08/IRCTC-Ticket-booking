package booking

import (
	"better-uptime/common/middleware"
	"better-uptime/common/routes"
	"better-uptime/config"
	db "better-uptime/internal/db/sqlc"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	config *config.Config
	store  db.Store
	Redis  redis.Client
}

func NewHandler(config *config.Config, store db.Store, Redis redis.Client) *Handler {
	return &Handler{
		config: config,
		store:  store,
		Redis:  Redis,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := routes.DefaultRouter()
	// without middleware

	router.Group(func(r chi.Router) {
		r.Use(middleware.TokenMiddleware(h.store))
		// with middleware

	})

	return router
}
