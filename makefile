.PHONY: build run clean up down kill rm ps network setup
.PHONY: db-create-schema dbs db-reset-fixtures db-reset-fixture dbf

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

kill: down
down:
	@docker-compose kill

rm:
	@$(MAKE) down
	@docker-compose rm

ps:
	@docker-compose ps -a

network:
	@if ! docker network ls | grep -q "binbogami"; then \
		@docker network create binbogami; \
	fi

setup:
	@cp -f .env.example .env
	@go get ./...
	@sudo -v
	@if ! grep -q "binbogami" /etc/hosts; then \
		sudo -- sh -c "echo '127.0.0.1	binbogami' >> /etc/hosts"; \
		sudo -- sh -c "echo '127.0.0.1	mariadb' >> /etc/hosts"; \
		sudo -- sh -c "echo '127.0.0.1	mailcatcher' >> /etc/hosts"; \
	fi
	@if ! docker network ls | grep -q "binbogami"; then \
  		@docker network create binbogami; \
  	fi

dbs: db-create-schema
db-create-schema:
	@mysql -u binbogami -plsktdngqcgnz3svj -h mariadb -P 3306 binbogami < var/resource/binbogami.sql

db-reset-fixture dbf: db-reset-fixtures
db-reset-fixtures:
	@mysql -u binbogami -plsktdngqcgnz3svj -h mariadb -P 3306 binbogami < var/resource/fixture.sql
