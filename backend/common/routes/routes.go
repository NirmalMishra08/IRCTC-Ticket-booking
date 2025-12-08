package routes

import (
	"better-uptime/common/logger"
	"fmt"

	"better-uptime/common/util"

	"net/http"

	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func DefaultRouter(middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.Heartbeat("/ping"),
		loggerMiddleware,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", util.VIN_HEADER},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)

	for _, middleware := range middlewares {
		router.Use(middleware)
	}

	return router
}

func loggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()

		defer func() {
			message := fmt.Sprintf("%s %s://%s%s %s - %d %dB in %v",
				r.Method, scheme, r.Host, r.RequestURI, r.Proto,
				ww.Status(), ww.BytesWritten(), time.Since(t1),
			)
			logger.Info(message)
		}()

		h.ServeHTTP(ww, r)
	})
}
