.PHONY: build
build:
	@go build -o ./bin/binbogami ./cmd/binbogami

.PHONY: run
run:
	@$(MAKE) build
	@bin/binbogami

.PHONY: up
up:
	@docker-compose up -d

.PHONY: down
down:
	@docker-compose kill

.PHONY: rm
rm:
	@$(MAKE) down
	@docker-compose rm

.PHONY: ps
ps:
	@docker-compose ps -a