# rapid-go

## motivation

rapid-go is a boilerplate that accelerates API development based on layered architecture and clarifying responsibilities.

## what is this

```
the boilerplate for monorepo application (support only http protocol)
```

- Base project is https://github.com/golang-standards/project-layout

## Apps

| Package                             | Localhost             | Prodction  |
| :---------------------------------- | :-------------------- | :--------- |
| **[[REST] default](./cmd/default)** | http://localhost:8080 | default.\* |

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

- local

```bash
> realize start
```

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
> sqlboiler -c ./db/authentication/sqlboiler.toml -o ./pkg/dbmodels/authentication psql
``` -->

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

### service layer

- write application logic using repository

### handler

- write the process about request and response
