package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	// "server/helpers"
	"server/helpers"
	"server/kafka"
	"server/models"

	// "strconv"

	"github.com/IBM/sarama"
	// "github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func SendMessage(producer sarama.SyncProducer, users []models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body := struct {
			FromID  int    `json:"from_id"`
			ToID    int    `json:"to_id"`
			Message string `json:"message"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			log.Error().Err(err).Msg("Error decoding message")
			helpers.ErrorJSON(w, errors.New("Error decoding message"), http.StatusBadRequest)
			return
		}

		log.Info().Msgf("Body: %v", body)

		err = kafka.SendKafkaMessage(
			producer,
			users, models.Notification{
				FromID:  body.FromID,
				ToID:    body.ToID,
				Message: body.Message,
			})

		if err != nil {
			log.Error().Err(err).Msg("Error sending message")
			helpers.ErrorJSON(w, errors.New("Error sending message"), http.StatusBadRequest)
			return
		}

		_ = helpers.WriteJSON(w, http.StatusOK, body)
	}
}
