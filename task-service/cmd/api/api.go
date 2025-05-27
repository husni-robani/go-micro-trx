package api

import (
	"errors"
	"net/http"
	"task-service/cmd/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type APIHandler struct {
	App *config.Config
}

// Routes
func (h *APIHandler) Routes() http.Handler{
	mux := chi.NewRouter()

	mux.Use(cors.Handler(
		cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders: []string{"Link"},
			AllowCredentials: true,
			MaxAge: 300,
		},
	))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/all", h.getAllTask)

	return mux
}

// HANDLERS
func (h *APIHandler) getAllTask(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.App.Models.Task.GetAll()
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, errors.New("failed to get tasks"))
		return
	}

	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "Get all task successful"
	responsePayload.Data = tasks

	h.writeResponse(w, http.StatusAccepted, responsePayload)
}