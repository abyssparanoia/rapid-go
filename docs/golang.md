# golang 環境

## environment (using direnv)
  - service account for gcp.
    - save it as serviceAccount.json

```bash
> cp .envrc.tepl .envrc
> direnv allow
```

## run

```bash
> make http.dev
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