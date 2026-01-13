.PHONY: service-up service-down database-up wait-db init-db migrations-up migrations-down

ENV_FILE := .env.development
DB_CONTAINER := postgres-dev

# Carrega variáveis do .env.development
ifneq (,$(wildcard $(ENV_FILE)))
	include $(ENV_FILE)
	export
endif

DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

service-up: database-up wait-db init-db migrations-up

service-down:
	docker compose -f ./infra/compose.yaml --env-file $(ENV_FILE) down

database-up:
	docker compose -f ./infra/compose.yaml --env-file $(ENV_FILE) up -d

wait-db:
	@echo -n "⏳ Waiting for database"
	@until docker exec $(DB_CONTAINER) pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) > /dev/null 2>&1; do \
		printf "."; \
		sleep 1; \
	done
	@echo ""
	@echo "✅ Database is ready"

init-db:
	@docker exec $(DB_CONTAINER) \
		psql -U $(POSTGRES_USER) -tc "SELECT 1 FROM pg_database WHERE datname='$(POSTGRES_DB)'" | grep -q 1 \
		|| docker exec $(DB_CONTAINER) psql -U $(POSTGRES_USER) -c "CREATE DATABASE $(POSTGRES_DB);"
	@echo "📦 Database $(POSTGRES_DB) ensured"

migrations-up:
	@echo "🚀 Running migrations..."
	goose -dir sql/schema postgres "$(DATABASE_URL)" up

migrations-down:
	@echo "⏬ Rolling back migrations..."
	goose -dir sql/schema postgres "$(DATABASE_URL)" down