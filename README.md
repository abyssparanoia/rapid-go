# rapid-go

## motivation

rapid-go is a boilerplate that accelerates API development based on layered architecture and clarifying responsibilities.

## what is this

```
the boilerplate for monorepo application (support only http protocol)
```

- Base project is https://github.com/golang-standards/project-layout

## Apps

| Package                                                            | Localhost              | Prodction            |
| :----------------------------------------------------------------- | :--------------------- | :------------------- |
| **[[REST] default](./cmd/default)**                                | http://localhost:8080  | default.\*           |
| **[[gRPC] default-grpc](./cmd/push-notification)**                 | http://localhost:50051 | default-grpc.\*      |
| **[[REST] default-grpc-json-transcoder](./cmd/push-notification)** | http://localhost:51051 | default-grpc-rest.\* |
| **[[REST] push-notification](./cmd/push-notification)**            | http://localhost:8081  | push-notification.\* |

## development

### Preparation

<!--
- generate rsa pem file

```bash
> openssl genrsa -out ./secret/catharsis-gcp.rsa 1024
> openssl rsa -in ./secret/catharsis-gcp.rsa  -pubout > ./secret/catharsis-gcp.rsa.pub
``` -->

- environment (using dotenv)
  - you should fix a host to default-db if you use docker-compose as server runtime

```bash
> cp .tmpl.env.default .env.default
```

### server starting

- docker

```bash
# build image
> docker-compose build

# container start
> docker-compose up -d
```

- example of default server

```bash
> curl --request GET 'http://localhost:8080/ping'
```

<!-- ### database

- generate server code by sql boiler

```bash
> make sqlboiler
``` -->

### testing

```bash
> docker-compose run --rm default-grpc-server ash -c "source .envrc && make test"
```

## production

### build

```bash
> docker build -f ./docker/production/default/Dockerfile .
```

## about layer

### infrastructure

- data layer
- It is responsibility to handle the data
- interested in database etc.

#### entity

- struct for setting the result of SQL etc....

#### infra/repository

- write the actual data manipulation process

### domain

#### model

- domain model

#### domain/repository

- write interface for infrastructure/repository and convert entity to domain

#### domain/service

- write application logic using repository

### usecase layer

- write usecase using repository and service

### handler

- write the process about request and response

### internal/pkg

- shared code
