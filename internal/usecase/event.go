package usecase

import (
	"encoding/json"

	"github.com/kiriksik/EventTracker/internal/domain/entity"
	"github.com/kiriksik/EventTracker/internal/domain/repository"
	"github.com/kiriksik/EventTracker/internal/infrastructure/kafka"
)

type EventUseCase struct {
	repo     repository.EventRepository
	producer *kafka.EventPublisher
}

func NewEventUseCase(repo repository.EventRepository, producer *kafka.EventPublisher) *EventUseCase {
	return &EventUseCase{repo: repo, producer: producer}
}

func (uc *EventUseCase) HandleEvent(event *entity.Event) error {
	// Сохранение в БД
	if err := uc.repo.Save(event); err != nil {
		return err
	}

	// Отправка в Кафку
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return uc.producer.Publish(event.ID, data)
}
