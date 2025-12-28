package api

import (
	"better-uptime/common/routes"

	"github.com/go-chi/chi/v5"
)

func (app *Server) routes() *chi.Mux {
	router := routes.DefaultRouter()

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/auth", app.authHandler.Routes())
		r.Mount("/train",app.trainHandler.Routes())
		r.Mount("/booking",app.bookingHandler.Routes())
	})

	return router
}
