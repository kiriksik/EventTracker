package clickhouse

import (
	"context"
	"database/sql"

	"github.com/kiriksik/EventTracker/internal/domain/entity"
	"github.com/kiriksik/EventTracker/internal/domain/repository"
	"go.uber.org/zap"
)

type eventRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewEventRepository(db *sql.DB, logger *zap.Logger) repository.EventRepository {
	return &eventRepo{db: db, logger: logger}
}

func (r *eventRepo) Save(event *entity.Event) error {
	_, err := r.db.ExecContext(context.Background(),
		"INSERT INTO events (id, type, payload, timestamp) VALUES (?, ?, ?, ?)",
		event.ID, event.Type, event.Payload, event.Timestamp,
	)
	if err != nil {
		r.logger.Error("failed to insert event", zap.Error(err))
	}
	return err
}
