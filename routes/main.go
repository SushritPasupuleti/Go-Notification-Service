package routes

import (
	// "encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	// "context"
	// "server/env"
	"server/handlers"
	"server/kafka"
	"server/models"

	// "strconv"

	"github.com/rs/zerolog/log"
)

func Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	users := []models.User{
		{ID: 1, Name: "Batman"},
		{ID: 2, Name: "Iron Man"},
		{ID: 3, Name: "Spiderman"},
		{ID: 4, Name: "Deadpool"},
	}

	store := &models.NotificationStore{
		Data: make(models.UserNotifications),
	}

	go kafka.SetupConsumerGroup(store)

	producer, err := kafka.SetupProducer()
	if err != nil {
		log.Fatal().Err(err).Msg("Error setting up Kafka producer")
	}
	// defer producer.Close()
	log.Info().Msg("Kafka Producer 📨 running")

	router.Group(func(r chi.Router) {
		router.Route("/producer", func(r chi.Router) {
			r.Post("/send", handlers.SendMessage(producer, users))
		})
	})

	router.Group(func(r chi.Router) {
		router.Route("/consumer", func(r chi.Router) {
			r.Get("/notifications/{id}", handlers.HandleNotification(store))
		})
	})

	return router
}
