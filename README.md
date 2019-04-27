# rapid-go

## motivation

rapid-go is a boilerplate that accelerates API development based on layered architecture and clarifying responsibilities.

## stack

- golang 1.11 (I will actively raise go version)
- mysql (correspondence such as firestore is easy)
- Chi (as Router)
- squirrel (as query builder)
- sqlx (map the result of sql to an object)
- gin (for hot reload ,not framwork)
- docker
- mockgen (generate mock codes from inteface)
- zap (as logger)
- firebase auth (as authenticate service)

## development

- init

```bash
make init
```

- build

```bash
> make build
```

- start

```bash
> make start
> curl http://localhost:3001/ping
```

- stop

```bash
> make down
```

- generate mock from interface (service,domain/repository)

```bash
> make mockgen_task
```

- run test

```bash
> make test
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
