package ws

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kiriksik/EventTracker/internal/telemetry"
	"go.uber.org/zap"
)

type Handler struct {
	eventUC EventUseCase
	metrics *telemetry.Metrics
	logger  *zap.Logger
}

type EventUseCase interface {
	ProcessEvent(ctx context.Context, name string, timestamp int64, payload map[string]interface{}) error
}

func NewHandler(eventUC EventUseCase, metrics *telemetry.Metrics, logger *zap.Logger) *Handler {
	return &Handler{
		eventUC: eventUC,
		metrics: metrics,
		logger:  logger,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IncomingEvent struct {
	Name      string                 `json:"name"`
	Timestamp int64                  `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade failed", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("websocket").Inc()
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.logger.Warn("Read error", zap.Error(err))
			h.metrics.ErrorsTotal.WithLabelValues("websocket").Inc()
			break
		}

		var event IncomingEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			h.logger.Warn("Invalid event JSON", zap.Error(err))
			h.metrics.ErrorsTotal.WithLabelValues("websocket").Inc()
			continue
		}

		err = h.eventUC.ProcessEvent(r.Context(), event.Name, event.Timestamp, event.Payload)
		if err != nil {
			h.logger.Error("Event processing failed", zap.Error(err))
			h.metrics.ErrorsTotal.WithLabelValues("usecase").Inc()
			continue
		}

		h.metrics.EventsTotal.WithLabelValues("ws").Inc()
	}
}
