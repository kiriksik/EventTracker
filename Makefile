.PHONY: help up down restart logs clean ps check-kafka clickhouse-cli redis-cli test check-env

# Цель по умолчанию (выводим справку)
help:
	@echo "Использование: make [цель]"
	@echo ""
	@echo "Цели:"
	@echo "  help           - Показать эту справку (цель по умолчанию)"
	@echp "  create-topic   - Создать топик events для Kafka"
	@echo "  up             - Запустить все сервисы (docker-compose up -d)"
	@echo "  down           - Остановить все сервисы (docker-compose down)"
	@echo "  restart        - Перезапустить сервисы (down + up)"
	@echo "  logs           - Показать логи всех сервисов"
	@echo "  clean          - Остановить сервисы и УДАЛИТЬ ВСЕ ДАННЫЕ (docker-compose down -v)"
	@echo "  ps             - Показать статус контейнеров"
	@echo "  check-kafka    - Проверить доступность Kafka (список топиков)"
	@echo "  clickhouse-cli - Зайти в консоль ClickHouse"
	@echo "  redis-cli      - Зайти в консоль Redis"
	@echo "  test           - Запустить тесты Go"
	@echo "  check-env      - Проверить наличие Docker и docker-compose"
	@echo ""
	@echo "Пример:"
	@echo "  make up"
	@echo "  make logs"

# Создать топик events в Kafka (если его нет)
create-topic:
	docker-compose exec kafka kafka-topics --bootstrap-server localhost:9092 --create --topic events --partitions 1 --replication-factor 1 || echo "Topic already exists or error"


# Запуск всех сервисов
up:
	docker-compose up -d
	sleep 10
	make create-topic

# Остановка всех сервисов
down:
	docker-compose down

# Перезапуск (остановка + запуск)
restart: down up

# Просмотр логов
logs:
	docker-compose logs -f

# Очистка ВСЕХ данных (включая volume'ы)
clean:
	docker-compose down -v
	rm -rf ./data  # Удаляем локальные данные, если есть

# Проверка состояния контейнеров
ps:
	docker-compose ps

# Проверка доступности Kafka
check-kafka:
	docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092

# Консоль ClickHouse
clickhouse-cli:
	docker-compose exec clickhouse clickhouse-client

# Консоль Redis
redis-cli:
	docker-compose exec redis redis-cli

# Запуск тестов
test:
	go test -v ./...

# Проверка окружения
check-env:
	@which docker || (echo "Error: Docker not found" && exit 1)
	@which docker-compose || (echo "Error: docker-compose not found" && exit 1)
	@echo "Dependencies OK"