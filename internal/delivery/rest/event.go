package rest

import (
	"encoding/json"
	"net/http"

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

	var event struct {
		Type      string                 `json:"type"`
		Timestamp int64                  `json:"timestamp"`
		Payload   map[string]interface{} `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := h.usecase.ProcessEvent(r.Context(), event.Type, event.Timestamp, event.Payload); err != nil {
		http.Error(w, "failed to process event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
