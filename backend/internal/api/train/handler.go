package train

import (
	"better-uptime/common/middleware"
	"better-uptime/common/routes"
	"better-uptime/config"
	db "better-uptime/internal/db/sqlc"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store  db.Store
	config *config.Config
}

func NewHandler(config *config.Config, store db.Store) *Handler {
	return &Handler{
		config: config,
		store:  store,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := routes.DefaultRouter()
	// without middleware

	router.Group(func(r chi.Router) {
		r.Use(middleware.TokenMiddleware(h.store))
		r.Get("/all-train", h.GetAllTrain)
		r.Post("/create-train", h.CreateTrain)
		r.Post("/coach-seat", h.CreateCoachesAndSeats)

	})

	return router
}
