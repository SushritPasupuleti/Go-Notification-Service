package handlers

import (
	// "encoding/json"
	"net/http"

	// "server/helpers"
	// "server/kafka"
	"server/helpers"
	"server/models"

	"github.com/go-chi/chi/v5"
)

func HandleNotification(store *models.NotificationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		notifications := store.GetNotifications(id)

		_ = helpers.WriteJSON(w, http.StatusOK, notifications)
	}
}
