package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/kiriksik/EventTracker/internal/configs"
	"github.com/kiriksik/EventTracker/internal/delivery/rest"
	"github.com/kiriksik/EventTracker/internal/delivery/ws"
	"github.com/kiriksik/EventTracker/internal/infrastructure/clickhouse"
	"github.com/kiriksik/EventTracker/internal/infrastructure/kafka"
	"github.com/kiriksik/EventTracker/internal/telemetry"
	"github.com/kiriksik/EventTracker/internal/usecase"
	"go.uber.org/zap"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := configs.LoadConfig("internal/configs/app.yaml")
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// 1. Метрики
	metrics := telemetry.NewMetrics()
	metrics.Register()
	telemetry.StartMetricsServer(":9090")

	// 3. Подключение к Clickhouse
	db, err := sql.Open("clickhouse",
		fmt.Sprintf("tcp://%s:%v?database=%s",
			cfg.ClickHouse.Host,
			cfg.ClickHouse.Port,
			cfg.ClickHouse.Database))
	if err != nil {
		logger.Fatal("failed to connect ClickHouse", zap.Error(err))
	}
	defer db.Close()

	// 4. Подключение к Kafka
	kafkaPublisher := kafka.NewEventPublisher(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer kafkaPublisher.Close()

	eventRepo := clickhouse.NewEventRepository(db, logger)
	eventUC := usecase.NewEventUseCase(eventRepo, kafkaPublisher, logger)

	// 5. Создание REST API и WebSocket
	handler := rest.NewEventHandler(eventUC, logger)
	mux := http.NewServeMux()
	wsHandler := ws.NewHandler(eventUC, metrics, logger)
	mux.HandleFunc("/ws", wsHandler.Handle)
	handler.RegisterRoutes(mux)

	srv := &http.Server{Addr: fmt.Sprintf(":%v", cfg.HTTP.Port), Handler: mux}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		logger.Info(fmt.Sprintf("starting server on :%v", cfg.HTTP.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", zap.Error(err))
	}
}
