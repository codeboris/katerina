package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/codeboris/katerina/internal/adapters/in/http/handlers"
	"github.com/codeboris/katerina/internal/adapters/in/http/middleware"
	"github.com/codeboris/katerina/internal/core/application/usecases"
)

func NewRouter(
	authUC *usecases.AuthUseCase,
	processUC *usecases.ProcessVoiceUseCase,
	audioUC *usecases.AudioUseCase,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Recoverer)
	r.Use(middleware.Logging)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	authHandler := handlers.NewAuthHandler(authUC)
	voiceHandler := handlers.NewVoiceHandler(processUC)
	audioHandler := handlers.NewAudioHandler(audioUC)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authUC))
			r.Post("/voice/process", voiceHandler.Process)
			r.Get("/audio/{id}", audioHandler.GetAudio)
		})
	})

	return r
}
