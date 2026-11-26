ENV_FILE ?= .env.local
-include $(ENV_FILE)

LOGGER_FOLDER ?= out/logs
export OTEL_EXPORTER_OTLP_ENDPOINT ?= http://127.0.0.1:4317

.PHONY: infra-build infra-up infra-down run-api run-ingest run-outbox run-normalizer run-embedding run-ocr run-finalizer run-web web-install web-build web-lint web-test logs-cleanup migrate-create migrate-up migrate-down migrate-force

run-web:
	npm --prefix web run dev

.PHONY: swagger
swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go --parseInternal -o docs

infra-build:
	@docker compose --env-file "$(ENV_FILE)" build

infra-up:
	@docker compose --env-file "$(ENV_FILE)" up -d

infra-down:
	@docker compose --env-file "$(ENV_FILE)" down

run-api:
	@set -ea; . "$(abspath $(ENV_FILE))"; LOGGER_FOLDER="$(abspath $(LOGGER_FOLDER))" go run ./cmd/api

run-ingest:
	@set -ea; . "$(abspath $(ENV_FILE))"; LOGGER_FOLDER="$(abspath $(LOGGER_FOLDER))" go run ./cmd/ingest

run-outbox:
	@set -ea; . "$(abspath $(ENV_FILE))"; LOGGER_FOLDER="$(abspath $(LOGGER_FOLDER))" go run ./cmd/outbox

run-normalizer:
	@set -ea; . "$(abspath $(ENV_FILE))"; LOGGER_FOLDER="$(abspath $(LOGGER_FOLDER))" go run ./cmd/normalizer

run-embedding:
	@set -ea; . "$(abspath $(ENV_FILE))"; LOGGER_FOLDER="$(abspath $(LOGGER_FOLDER))" go run ./cmd/embedding

run-ocr:
	@set -ea; . "$(abspath $(ENV_FILE))"; LOGGER_FOLDER="$(abspath $(LOGGER_FOLDER))" go run ./cmd/ocr

run-finalizer:
	@set -ea; . "$(abspath $(ENV_FILE))"; LOGGER_FOLDER="$(abspath $(LOGGER_FOLDER))" go run ./cmd/finalizer

logs-cleanup:
	@read -p "Clear all logs? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ./out/logs && \
		echo "Done"; \
	else \
		echo "Canceled"; \
	fi;

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Parameter seq not found"; \
		exit 1; \
	fi
	@docker run --rm \
		-v "$(abspath migrations):/migrations" \
		migrate/migrate \
		create -ext sql -dir /migrations -seq $(seq)

migrate-up:
	@docker run --rm \
		--network image_uploader_net \
		-v "$(abspath migrations):/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}" \
		up

migrate-down:
	@docker run --rm \
		--network image_uploader_net \
		-v "$(abspath migrations):/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}" \
		down 1

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Parameter version not found"; \
		exit 1; \
	fi
	@docker run --rm \
		--network image_uploader_net \
		-v "$(abspath migrations):/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}" \
		force $(version)
