# =========================
# ENV
# =========================
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

APP_NAME    ?= url_shortener
APP_VERSION ?= 0.0.1

BIN_DIR ?= bin
CMD_DIR ?= ./cmd/app

TARGET_BIN = $(BIN_DIR)/$(APP_NAME)-$(APP_VERSION)

DB_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB_NAME)?sslmode=disable

# =========================
# DOCKER COMPOSE CORE
# =========================

.PHONY: up
up:
	docker compose up --build

.PHONY: down
down:
	docker compose down

.PHONY: reset
reset:
	docker compose down -v
	docker compose up --build

.PHONY: logs
logs:
	docker compose logs -f

.PHONY: ps
ps:
	docker compose ps

# =========================
# DATABASE (ONLY)
# =========================

.PHONY: db-up
db-up:
	docker compose up -d postgres

.PHONY: db-down
db-down:
	docker compose stop postgres


# =========================
# REDIS (ONLY)
# =========================

.PHONY: redis-up
redis-up:
	docker compose up -d redis

.PHONY: redis-down
redis-down:
	docker compose stop redis

.PHONY: redis-restart
redis-restart:
	docker compose restart redis

.PHONY: redis-logs
redis-logs:
	docker compose logs -f redis

.PHONY: redis-cli
redis-cli:
	docker exec -it redis redis-cli -a $(REDIS_PASSWORD)


# =========================
# MIGRATIONS
# =========================

.PHONY: migrate-up
migrate-up:
	docker compose run --rm migrate \
		-path=/migrations \
		-database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	docker compose run --rm migrate \
		-path=/migrations \
		-database "$(DB_URL)" down 1

.PHONY: migrate-reset
migrate-reset:
	docker compose run --rm migrate \
		-path=/migrations \
		-database "$(DB_URL)" drop -f

.PHONY: migrate-version
migrate-version:
	docker compose run --rm migrate \
		-path=/migrations \
		-database "$(DB_URL)" version

.PHONY: migrate-force
migrate-force:
	docker compose run --rm migrate \
		-path=/migrations \
		-database "$(DB_URL)" force $(VERSION)

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

# =========================
# BUILD (LOCAL)
# =========================

.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	go build -o $(TARGET_BIN) $(CMD_DIR)/main.go

.PHONY: run
run: build
	./$(TARGET_BIN)

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

# =========================
# DEV WORKFLOW
# =========================

.PHONY: dev
dev:
	docker compose up $(if $(BUILD), --build,) $(if $(DETACH), -d,)

.PHONY: dev-reset
dev-reset:
	docker compose down -v
	docker compose up --build