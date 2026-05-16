# Agents

## Project Context
- **Language**: Go 1.25+
- **Frameworks**: Labstack Echo v5
- **Primary database**: MariaDB 12.1+

## Project Structure
- `bin/`: Service binaries
- `cmd/`: Main files for the services
- `database/`: SQL queries for database schema and fixtures
- `internal/`: Core code
- `templates/`: HTML templates
- `var/`: Runtime generated code such as logs and Docker volumes

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
    - `FindBy*`: Find single entity by something
    - `FindManyBy*`: Find many entities by something
    - `Save`: Persist entity (create or update)
    - `Create(entity *Entity) error`: Persist new entity in the database
    - `Update(entity *Entity) error`: Persist existing entity in the database
    - `Find(id uuid.UUID) (*Entity, Error)`: Find single entity by UUID
  - **REST API Handlers**:
    - **GET** for multiple entities assigned endpoint `/resources` and function name `index`
    - **GET** for sinle entity assigned endpoint `/resources/{:id}` and function name `show`
    - **POST** for creation assigned endpoint `/resources` and function name `create`
    - **PUT** for ammendment assigned endpoint `/resources/{:id}` and function name `update`
    - **DELETE** assigned endpoint `/resources/{:id}` and function name `destroy`

## Agent Instructions
- When writing tests, place them in `_test.go` files in the same package.
