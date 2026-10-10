# Binbogami

A kami or god who inhabits a human being or their house to bring misery and poverty.

Personal finance management tool which should replace existing Google document.

## Architecture & Services

Binbogami is designed with a service-oriented architecture (SOA). Each service has its own dedicated entrypoint in `cmd/` and runs on a configured port, or can be run together via the primary combined service:

- **`binbogami`**: Primary API gateway / combined service (port `8101`).
- **`auth`**: Authentication, user invitation, and password reset service (port `8101`).
- **`finance`**: Finance management service for books, categories, locations, and transaction entries (port `8104`).
- **`user`**: User profile management and configuration service (port `8103`).
- **`worker`**: Redis-backed [Asynq](https://github.com/hibiken/asynq) background task worker for OCR receipt processing ([Tesseract](https://github.com/otiai10/gosseract)) and AI text extraction ([Google Gemini](https://github.com/google/generative-ai-go)).
- **`content`**: Content service (stub).

### Supporting Infrastructure

Local infrastructure services run in Docker via [compose.yaml](file:///Users/jurgis/Develop/binbogami/compose.yaml):

- **MariaDB 12.3**: Primary relational database (port `3306`).
- **Redis 8**: In-memory broker for asynchronous background workers (port `6379`).
- **SeaweedFS**: S3-compatible distributed object storage for receipt images and file uploads (ports `8333`, `8888`, `9333`).
- **Grafana Loki 3.7 & Grafana 13.0**: Centralized logging system via custom Go `slog` handler and dashboard visualization (ports `3100`, `3000`).
- **Mailcatcher**: Local SMTP server (port `1025`) and web interface (port `1080`) for email delivery testing.

## User Interfaces

- [Web application](https://github.com/jurgisjaska/binbogami-web) built with Vue 3.
- [Desktop application](https://github.com/jurgisjaska/binbogami-desktop) built with GTK 4.
- Mobile (Android) application built with Kotlin.

## Project Structure

```
├── bin/          # Compiled service binaries
├── cmd/          # Service entrypoints (auth, binbogami, content, finance, user, worker)
├── database/     # SQL schema (schema.sql) and test fixtures (fixtures.sql)
├── internal/     # Core application packages (api, database, handlers, queue, service)
├── templates/    # HTML email templates (invitation.html, reset_password.html)
└── var/          # Runtime generated data (logs, Docker container volumes)
```

## Getting Started

### Prerequisites

- [Go](https://go.dev/) 1.25+
- [Docker](https://www.docker.com/) and Docker Compose
- MySQL/MariaDB client (`mysql` CLI to apply schema and fixtures)
- (Optional, for `worker` OCR): Tesseract and Leptonica development libraries (e.g. `brew install tesseract leptonica`)

### Initial Setup

Run the automated setup command to initialize your local environment:

```shell
make setup
```

This command performs the following initial steps:
1. Copies `.env.example` to `.env`.
2. Installs Go dependencies (`go get ./...`).
3. Appends local host entries (`binbogami`, `mariadb`, `mailcatcher`, `seaweedfs`) to `/etc/hosts` pointing to `127.0.0.1` (requires `sudo`).
4. Creates the `binbogami` external Docker network if it does not already exist.

### Start Supporting Containers

Start the Docker containers:

```shell
make up
```

To stop containers, run:

```shell
make down
```

### Initialize Database

Apply the database schema and load test fixtures:

```shell
make schema
make fixtures
```

### Build and Run

#### Running Services Locally

Run the default combined service (`binbogami`):

```shell
make run
```

Once running, the API can be accessed at [http://localhost:8101](http://localhost:8101/).

To run specific microservices, pass the `SERVICE` variable:

```shell
# Run a single service
make run SERVICE=auth

# Run multiple services simultaneously
make run SERVICE=auth,finance,user,worker
```

#### Building Binaries

Build service binaries into the `bin/` directory:

```shell
# Build default binary (binbogami)
make build

# Build specific service(s)
make build SERVICE=auth
make build SERVICE=auth,finance,worker
```

All compiled binaries are output to the `bin/` directory.

### Testing and Code Quality

Run tests with coverage reporting:

```shell
make test
```

Or run package-focused tests:

```shell
go test -v ./internal/<package>/...
```

Format code and tidy dependencies:

```shell
go fmt ./...
go mod tidy
```

## Code Style and Naming Conventions

### Go Standards

- Use `camelCase` for internal and `PascalCase` for exported members.
- Follow Go acronym rules (e.g., `JSONData`, `UUID`, not `JsonData`, `Uuid`).
- Use `CreateXxx` for constructor functions (e.g., `CreateBook`, `CreateUser`, `CreateConfig`).

### REST API Handlers

Route parameters follow [Labstack Echo](https://echo.labstack.com/) syntax (`:id`):

| Method | Endpoint        | Handler Function | Description              |
|--------|-----------------|------------------|--------------------------|
| GET    | `/resources`     | `index`          | Fetch multiple entities  |
| GET    | `/resources/:id` | `show`           | Fetch single entity      |
| POST   | `/resources`     | `create`         | Create new entity        |
| PUT    | `/resources/:id` | `update`         | Amend existing entity    |
| DELETE | `/resources/:id` | `destroy`        | Delete entity            |

### JSON API

API payloads follow `snake_case` naming conventions:

```json
{
  "description": "this is a description",
  "user_id": "3f89f6b5-3760-4d85-8b2e-31cca32e4913"
}
```

### Database

- Database tables use plural names (e.g., `users`, `books`, `categories`, `entries`).
- Column names use `snake_case` with suffixes for timestamps (e.g., `created_at`, `updated_at`, `deleted_at`).

### Repositories

Conventions for repository method naming and return types:

| Method | Return | Description |
|---|---|---|
| `Find(id uuid.UUID)` | `(*Entity, error)` | Find a single entity by UUID |
| `FindBy*({value})` | `(*Entity, error)` | Find a single entity by attribute |
| `FindManyBy*({value})` | `(Collection, error)` | Find many entities by attribute |
| `Create(entity *Entity)` or `Create(model *models.Entity)` | `error` or `(*Entity, error)` | Persist **new** entity in the database |
| `Update(entity *Entity)` | `error` | Persist **existing** entity in the database |

- Repository interfaces are named using pattern `*Repository` (e.g., `UserRepository`, `BookRepository`).
- Repositories are structs named `Repository` if the package contains only one repository.
- If the package contains multiple repositories, the name of the repository is formed using pattern `EntityRepository` and interfaces include a package name.
