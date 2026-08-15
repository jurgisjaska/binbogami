# Binbogami

A kami or god who inhabits a human being or their house to bring misery and poverty.

Personal finance management tool which should replace existing Google document.

## Installation

## Development

Configure the environment by creating a copy the distributed environment configuration

```shell
cp .env.example .env
```

Create `binbogami` Docker network

```shell
docker network create binbogami
```

Network existance can be checked using command `docker network ls | grep binbogami`

Start [Docker](https://www.docker.com) containers for database and other services

```shell
make up
```

Build the application and start web server

```shell
make run
```

After the web server is up and running the API can be reached by going to [localhost:8101](http://localhost:8101/)

### Setup

There is a single command for initial setup

```shell
make setup
```

Which performs all the steps needed for the first time use.
After it's execution continue with the starting of the Docker containers

## Code Style and Naming Conventions

### Handler Methods

Recommended handler method names for HTTP methods.

| Method | Endpoint        | Name    |
|--------|-----------------|---------|
| GET    | /resources      | index   |
| GET    | /resources/{id} | show    |
| POST   | /resources      | create  |
| PUT    | /resources      | update  |
| DELETE | /resources/{id} | destroy |

### JSON API

Service communicates using JSON API where JSON is formated to follow `camelCase` naming.

```json
{
  "description": "this is a description",
  "userId": "3f89f6b5-3760-4d85-8b2e-31cca32e4913"
}
```

### Database

Database uses plural table names (ex.: `users`) and snake_case names with suffixes for fields (ex.: `created_at`)

### Repositories

Convention for repository method naming and return types.

| Method | Return | Description |
|---|---|---|
| `Find(id uuid.UUID)` | `(*Entity, error)` | Find a single entity by UUID |
| `FindBy*({value})` | `(*Entity, error)` | Find a single entity by something |
| `FindManyBy*({value})` | `(Collection, error)` | Find many entities by something |
| `Save({entity})` | `error` | Persist entity (create or update) |
| `Create(entity *Entity)` | `error` | Persist **new** entity in the database |
| `Update(entity *Entity)` | `error` | Persist **existing** entity in the database |

- Repository interfaces are named using pattern `*Repository`.
- Repositories are structs named `Repository` if the package contains only one repository.
- If the package contains multiple repositories, the name of the repository is formed using pattern `EntityRepository` and interfaces include a package name.

