package tatkal

import (
	"better-uptime/common/middleware"
	"better-uptime/common/ratelimiter"
	"better-uptime/common/routes"
	"better-uptime/config"
	db "better-uptime/internal/db/sqlc"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	store       db.Store
	config      *config.Config
	RateLimiter *ratelimiter.RedisRateLimiter
	Redis       redis.Client
}

func NewHandler(config *config.Config, store db.Store, Redis redis.Client) *Handler {
	return &Handler{
		config:      config,
		store:       store,
		RateLimiter: ratelimiter.NewRedisRateLimiter(&Redis),
		Redis:  Redis,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := routes.DefaultRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.TokenMiddleware(h.store))

	})

	return router
}
