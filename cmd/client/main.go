package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kiriksik/EventTracker/internal/configs"
	"github.com/kiriksik/EventTracker/internal/infrastructure/kafka"
	"go.uber.org/zap"
)

func main() {
	// 1. Конфиг
	cfg, err := configs.LoadConsumerConfig("internal/configs/consumer.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Инициализация логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Starting Kafka consumer...", zap.String("service", cfg.ServiceName))

	// 3. Создание Kafka Consumer-а
	consumer := kafka.NewEventConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, logger)
	ctx, cancel := context.WithCancel(context.Background())

	// 4. Запуск Consumer-а
	go func() {
		if err := consumer.Consume(ctx); err != nil {
			logger.Error("error consuming messages", zap.Error(err))
			cancel()
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	logger.Info("Shutting down consumer gracefully...")
	cancel()
}
