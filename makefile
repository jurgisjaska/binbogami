.PHONY: build run clean up down rm ps setup

NAME=binbogami

build:
	@go build -o ./bin/${NAME} ./cmd/${NAME}

run:
	@$(MAKE) build
	@bin/${NAME}

test:
	@go test -v -cover -coverprofile=coverage.out ./...

clean:
	@rm -f bin/${NAME}

up:
	@docker-compose up -d

down:
	@docker-compose kill

rm:
	@$(MAKE) down
	@docker-compose rm

ps:
	@docker-compose ps -a

setup:
	@cp .env.example .env
	@go get ./...
	@docker network create binbogami