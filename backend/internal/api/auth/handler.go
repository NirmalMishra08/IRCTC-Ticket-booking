package auth

import (
	"better-uptime/common/middleware"
	"better-uptime/common/routes"
	"better-uptime/config"
	db "better-uptime/internal/db/sqlc"
	"go/token"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	config *config.Config
	store  db.Store
}


type HandlerConfig struct {
	Config     *config.Config
	Store      db.Store
	TokenMaker token.Token
}

func NewHandler(config *config.Config, store db.Store) *Handler {
	return &Handler{
		config: config,
		store:  store,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := routes.DefaultRouter()
	router.Post("/login", h.Login)

	router.Group(func(r chi.Router) {
		r.Use(middleware.TokenMiddleware(h.store))

	})

	return router
}
