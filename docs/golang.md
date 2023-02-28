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
# create new migration files in internal/infrastructure/database/migration/files
> make migrate.create

# migrate up
> make migrate.up
```

## generate sqlboiler

```bash
> make generate.sqlboiler
```

## create seed staff

- create to local database and remote firebase auth
  - If a staff has already been created in firebase, this cli create a staff only in the local database.

```bash
> make build
> ./.bin/app-cli task create-root-staff --email <email address> --password <passowrd>
```
