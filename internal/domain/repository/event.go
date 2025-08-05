package repository

import "github.com/kiriksik/EventTracker/internal/domain/entity"

type EventRepository interface {
	Save(event *entity.Event) error
}
