.PHONY: build run clean up down rm ps

NAME=binbogami

build:
	@go build -o ./bin/${NAME} ./cmd/${NAME}

run:
	@$(MAKE) build
	@bin/${NAME}

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