# Agents

## Project Context
- **Language**: Go 1.26+
- **Frameworks**: Labstack Echo v5
- **Primary database**: MariaDB 12.1+
- **Logging**: Grafana Loki

## Project Structure
- `bin/`: Service binaries
- `cmd/`: Main files for the services
- `database/`: SQL queries for database schema and fixtures
- `internal/`: Core code
- `templates/`: HTML templates
- `var/`: Runtime generated code such as logs and Docker volumes

## Service Oriented Architecture
- Every service is a separate `main.go` in the `cmd/` directory.
- Every service has its own port for API.

## UI
- [Web application](https://github.com/jurgisjaska/binbogami-web) built with Vue 3 
- [Desktop application](https://github.com/jurgisjaska/binbogami-desktop) built with GTK 4.
- Mobile (Android) application built with Kotlin.

## Development Workflow
- **Linting**: Always run `go fmt ./...` before suggesting code.
- **Dependencies**: Use `go mod tidy` after adding new imports.
- **Testing**: Use standard `go test -v -cover -coverprofile=coverage.out ./...`. Prefer table-driven tests.
- **Make**: Use GNU Make and commands in makefile for workflow automation.

## Coding Standards
- **Naming**:
  - Use camelCase for internal and PascalCase for exported members. Follow Go acronym rules (e.g., `JSONData`, not `JsonData`).
  - Use `CreateXxx` for constructor functions
  - **Repositories**:
    - `FindBy*`: Find a single entity by something
    - `FindManyBy*`: Find many entities by something
    - `Save`: Persist entity (create or update)
    - `Create(entity *Entity) error`: Persist new entity in the database
    - `Update(entity *Entity) error`: Persist existing entity in the database
    - `Find(id uuid.UUID) (*Entity, Error)`: Find a single entity by UUID
    - Repository interfaces are named using pattern `*Repository`
    - Repositories are structs named `Repository` if the package contains only one repository.
    - If the package contains multiple repositories, the name of the repository is formed using pattern `EntityRepository` and interfaces include a package name.
  - **REST API Handlers**:
    - **GET** for multiple entities assigned endpoint `/resources` and function name `index`
    - **GET** for sinle entity assigned endpoint `/resources/{:id}` and function name `show`
    - **POST** for creation assigned endpoint `/resources` and function name `create`
    - **PUT** for ammendment assigned endpoint `/resources/{:id}` and function name `update`
    - **DELETE** assigned endpoint `/resources/{:id}` and function name `destroy`

## Agent Instructions
- When writing tests, place them in `_test.go` files in the same package.
