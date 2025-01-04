package migration

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	constant_files "github.com/abyssparanoia/rapid-go/db/main/constants"
	migration_files "github.com/abyssparanoia/rapid-go/db/main/migrations"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger/logger_field"
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
	tables, err := databaseCli.DB.Query("SHOW TABLES") //nolint:rowserrcheck
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

func RunSyncConstants() {
	e := &environment.DatabaseEnvironment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}
	ctx := context.Background()

	l := logger.New()
	ctx = logger.ToContext(ctx, l)

	logger.L(ctx).Info("start sync constants")

	databaseCli := database.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, true)

	dirEntries, err := constant_files.EmbedConstants.ReadDir(".")
	if err != nil {
		panic(err)
	}

	tx, err := databaseCli.DB.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	var syncErr error
	defer func() {
		// If there's an error, rollback the transaction, else commit it
		if syncErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.L(ctx).Error("Failed to rollback transaction", logger_field.Error(rollbackErr))
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				logger.L(ctx).Error("Failed to commit transaction", logger_field.Error(commitErr))
			}
		}
	}()

	for _, entry := range dirEntries {
		if strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml") {
			syncErr = syncConstantByYaml(ctx, tx, entry)
			if syncErr != nil {
				return
			}
		}
	}
}
