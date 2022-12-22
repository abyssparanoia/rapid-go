package migration

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/caarlos0/env/v6"
	"github.com/pressly/goose/v3"
)

//go:embed files/*.sql
var embedMigrations embed.FS

func RunNewFile(fileName string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := goose.Create(nil, filepath.Join(dir, "internal/infrastructure/database/migration/files"), fileName, "sql"); err != nil {
		panic(err)
	}
}

func RunUp() {
	e := &environment.Environment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}

	logger := logger.New()

	logger.Info("start database schema migration")

	databaseCli := database.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase)

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		panic(err)
	}

	if err := goose.Up(databaseCli.DB, "files"); err != nil {
		panic(err)
	}

	logger.Info("complete database schema migration")
}
