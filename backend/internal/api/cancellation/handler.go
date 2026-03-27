package cancellation

import (
	"better-uptime/common/kafka"
	"better-uptime/common/middleware"
	"better-uptime/common/routes"
	"better-uptime/config"
	db "better-uptime/internal/db/sqlc"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store  db.Store
	config *config.Config
    Kafka  kafka.Producer
}

func NewHandler(config *config.Config, store db.Store, Kafka kafka.Producer) *Handler {
	return &Handler{
		config: config,
		store:  store,
		Kafka: Kafka,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := routes.DefaultRouter()
	// without middleware

	router.Group(func(r chi.Router) {
		r.Use(middleware.TokenMiddleware(h.store))
	})

	return router
}
