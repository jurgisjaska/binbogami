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

