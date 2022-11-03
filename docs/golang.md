# golang 環境

## environment (using direnv)

- service account for gcp.
  - save it as serviceAccount.json

```bash
> cp .envrc.tepl .envrc
> direnv allow
```

- replce GCP_PROJECT_ID and FIREBASE_CLIENT_KEY

## run

```bash
> docker-compose up -d # database
> make migrate.up
> make http.dev
> curl http://localhost:8080
```

## lint

```bash
> make lint.go

> make lint.proto
```

## generate protocol buffer + openapi(v2)

```bash
> make generate.buf
```

## db migration

```bash
> make migrate.up
```

## generate sqlboiler

```bash
> make generate.sqlboiler
```

## create seed user

- create to local database and remote firebase auth
  - If a user has already been created in firebase, this cli create a user only in the local database.

```bash
> make build
> ./.bin/app-cli task create-root-user --email <email address> --password <passowrd>
```
