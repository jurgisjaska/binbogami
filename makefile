.PHONY: build run up down rm ps

build:
	@go build -o ./bin/binbogami ./cmd/binbogami

run:
	@$(MAKE) build
	@bin/binbogami

up:
	@docker-compose up -d

down:
	@docker-compose kill

rm:
	@$(MAKE) down
	@docker-compose rm

ps:
	@docker-compose ps -a