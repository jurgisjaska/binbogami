.PHONY: build run clean up down kill rm ps network setup
.PHONY: schema fixtures

PROJECT=binbogami
SERVICE?=binbogami
include .env
export $(shell sed 's/=.*//' .env)

build:
	@go build -o ./bin/${SERVICE} ./cmd/${SERVICE}

run:
	@$(MAKE) build SERVICE=${SERVICE}
	@bin/${SERVICE}

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
	@if ! grep -q ${PROJECT} /etc/hosts; then \
		sudo -- sh -c "echo '127.0.0.1	${PROJECT}' >> /etc/hosts"; \
		sudo -- sh -c "echo '127.0.0.1	mariadb' >> /etc/hosts"; \
		sudo -- sh -c "echo '127.0.0.1	mailcatcher' >> /etc/hosts"; \
	fi
	@if ! docker network ls | grep -q ${PROJECT}; then \
  		docker network create ${PROJECT}; \
  	fi

schema:
	@mysql -u $(DATABASE_USERNAME) -p$(DATABASE_PASSWORD) -h $(DATABASE_HOSTNAME) -P $(DATABASE_PORT) $(DATABASE_NAME) < database/schema.sql

fixtures:
	@mysql -u $(DATABASE_USERNAME) -p$(DATABASE_PASSWORD) -h $(DATABASE_HOSTNAME) -P $(DATABASE_PORT) $(DATABASE_NAME) < database/fixtures.sql
