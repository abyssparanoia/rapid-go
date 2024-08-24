package migration

import (
	"os"
	"path/filepath"

	migration_files "github.com/abyssparanoia/rapid-go/db/main/migrations"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/caarlos0/env/v11"
	"github.com/pressly/goose/v3"
)

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
	e := &environment.DatabaseEnvironment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}

	logger := logger.New()

	logger.Info("start database schema migration")

	databaseCli := database.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, true)

	goose.SetBaseFS(migration_files.EmbedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		panic(err)
	}

	if err := goose.Up(databaseCli.DB, "."); err != nil {
		panic(err)
	}

	logger.Info("complete database schema migration")
}

func RunExtractSchema() {
	e := &environment.DatabaseEnvironment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}

	logger := logger.New()

	logger.Info("start extracting database schema")

	databaseCli := database.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, true)

	//nolint:execinquery
	tables, err := databaseCli.DB.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}
	defer tables.Close()

	file, err := os.Create("./db/main/schema.sql")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var tableName string
	for tables.Next() {
		if err := tables.Scan(&tableName); err != nil {
			panic(err)
		}
		var createStatement string
		//nolint:execinquery
		if err := databaseCli.DB.QueryRow("SHOW CREATE TABLE "+tableName).Scan(&tableName, &createStatement); err != nil {
			panic(err)
		}
		_, err := file.WriteString(createStatement + ";\n\n")
		if err != nil {
			panic(err)
		}
	}
}
