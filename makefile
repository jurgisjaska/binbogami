.PHONY: build run clean up down kill rm ps network setup
.PHONY: schema fixtures

PROJECT=binbogami
SERVICE?=binbogami
comma := ,
SERVICES := $(strip $(subst $(comma), ,$(SERVICE)))

include .env
export $(shell sed 's/=.*//' .env)

build:
	@for s in $(SERVICES); do \
		echo "Building $$s..."; \
		go build -o ./bin/$$s ./cmd/$$s || exit 1; \
	done

run:
	@$(MAKE) build SERVICE="$(SERVICE)"
	-@pids=""; \
	trap 'kill $$pids 2>/dev/null; trap - INT TERM EXIT; exit 0' INT TERM EXIT; \
	for s in $(SERVICES); do \
		echo "Starting $$s..."; \
		./bin/$$s & \
		pids="$$pids $$!"; \
	done; \
	wait 2>/dev/null || true

test:
	@go test -v -cover -coverprofile=coverage.out ./...

up:
	@docker-compose up -d

kill: down
down:
	@docker-compose kill

network:
	@if ! docker network ls | grep -q ${PROJECT}; then \
		docker network create ${PROJECT}; \
	fi

setup:
	@cp -f .env.example .env
	@go get ./...
	@sudo -v
	@for host in $(PROJECT) mariadb mailcatcher; do \
		if ! grep -v '^[[:space:]]*#' /etc/hosts | grep -E -q "[[:space:]]$$host([[:space:]]|$$)"; then \
			sudo -- sh -c "echo '127.0.0.1	$$host' >> /etc/hosts"; \
		fi; \
	done
	@if ! docker network ls | grep -q ${PROJECT}; then \
		docker network create ${PROJECT}; \
	fi

schema:
	@echo "Applying database schema..."
	@MYSQL_PWD="$(DATABASE_PASSWORD)" mysql -u $(DATABASE_USERNAME) -h $(DATABASE_HOSTNAME) -P $(DATABASE_PORT) $(DATABASE_NAME) < database/schema.sql
	@echo "Database schema applied."

fixtures:
	@echo "Loading database fixtures..."
	@MYSQL_PWD="$(DATABASE_PASSWORD)" mysql -u $(DATABASE_USERNAME) -h $(DATABASE_HOSTNAME) -P $(DATABASE_PORT) $(DATABASE_NAME) < database/fixtures.sql
	@echo "Database fixtures loaded."

