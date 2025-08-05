package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kiriksik/EventTracker/internal/domain/entity"
	"github.com/kiriksik/EventTracker/internal/usecase"
	"go.uber.org/zap"
)

type EventHandler struct {
	usecase *usecase.EventUseCase
	logger  *zap.Logger
}

func NewEventHandler(uc *usecase.EventUseCase, logger *zap.Logger) *EventHandler {
	return &EventHandler{usecase: uc, logger: logger}
}

func (h *EventHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/event", h.HandleIncomingEvent)
}

func (h *EventHandler) HandleIncomingEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	event := &entity.Event{
		ID:        uuid.New().String(),
		Type:      input.Type,
		Payload:   input.Payload,
		Timestamp: time.Now(),
	}

	if err := h.usecase.HandleEvent(event); err != nil {
		h.logger.Error("failed to handle event", zap.Error(err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
