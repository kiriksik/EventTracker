package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/kiriksik/EventTracker/internal/domain/entity"
	"github.com/kiriksik/EventTracker/internal/domain/repository"
	"github.com/kiriksik/EventTracker/internal/infrastructure/kafka"
	"go.uber.org/zap"
)

type EventUseCase struct {
	repo     repository.EventRepository
	producer *kafka.EventPublisher
	logger   *zap.Logger
}

func NewEventUseCase(repo repository.EventRepository, producer *kafka.EventPublisher, logger *zap.Logger) *EventUseCase {
	return &EventUseCase{repo: repo, producer: producer, logger: logger}
}

func (uc *EventUseCase) ProcessEvent(ctx context.Context, typ string, timestamp int64, payload map[string]interface{}) error {
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		uc.logger.Error("failed to marshal payload", zap.Error(err))
		return err
	}

	event := &entity.Event{
		ID:        uuid.NewString(),
		Type:      typ,
		Payload:   string(payloadStr),
		Timestamp: time.Unix(timestamp, 0),
	}

	// 1. Сохраняем в ClickHouse
	if err := uc.repo.Save(event); err != nil {
		uc.logger.Error("failed to save event to DB", zap.Error(err))
	}

	// 2. Публикуем в Kafka
	data, err := json.Marshal(event)
	if err != nil {
		uc.logger.Error("failed to marshal event for Kafka", zap.Error(err))
		return err
	}

	if err := uc.producer.Publish(ctx, event.ID, data); err != nil {
		uc.logger.Error("failed to publish event", zap.Error(err))
		return err
	}

	return nil
}
