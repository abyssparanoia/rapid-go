package migration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	constant_files "github.com/abyssparanoia/rapid-go/db/mysql/constants"
	migration_files "github.com/abyssparanoia/rapid-go/db/mysql/migrations"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger/logger_field"
	"github.com/caarlos0/env/v11"
	"github.com/pressly/goose/v3"
	"gopkg.in/yaml.v3"
)

func RunNewFile(fileName string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := goose.Create(nil, filepath.Join(dir, "internal/infrastructure/mysql/migration/files"), fileName, "sql"); err != nil {
		panic(err)
	}
}

func RunUp() {
	e := &environment.DatabaseEnvironment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}

	logger := logger.New(environment.MinLogLevelInfo)

	logger.Info("start database schema migration")

	databaseCli := mysql.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, true)

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

	logger := logger.New(environment.MinLogLevelInfo)

	logger.Info("start extracting database schema")

	databaseCli := mysql.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, true)

	//nolint:execinquery
	tables, err := databaseCli.DB.Query("SHOW TABLES") //nolint:rowserrcheck
	if err != nil {
		panic(err)
	}
	defer tables.Close()

	file, err := os.Create("./db/mysql/schema.sql")
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

// constantData represents the structure of a YAML file for a constants table.
type constantData []*struct {
	Table  string   `yaml:"table"`
	Values []string `yaml:"values"`
}

func RunSyncConstants() {
	e := &environment.DatabaseEnvironment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}
	ctx := context.Background()

	l := logger.New(environment.MinLogLevelInfo)
	ctx = logger.ToContext(ctx, l)

	databaseCli := mysql.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, true)
	runSyncConstantsWithContext(ctx, databaseCli)
}

func runSyncConstantsWithContext( //nolint:gocognit
	ctx context.Context,
	databaseCli *mysql.Client,
) {
	logger.L(ctx).Info("start sync constants")

	var data constantData
	if err := yaml.Unmarshal(constant_files.EmbedConstants, &data); err != nil {
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
			logger.L(ctx).Error("Failed to sync constants", logger_field.Error(syncErr))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.L(ctx).Error("Failed to rollback transaction", logger_field.Error(rollbackErr))
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				logger.L(ctx).Error("Failed to commit transaction", logger_field.Error(commitErr))
			}
		}
	}()

	if err := func() error {
		for _, dst := range data {
			// Fetch current IDs from the database
			query := fmt.Sprintf("SELECT id FROM %s", dst.Table) //nolint:gosec
			rows, err := tx.Query(query)                         //nolint:rowserrcheck
			if err != nil {
				return errors.InternalErr.Wrap(err)
			}
			defer rows.Close()

			currentIDs := make(map[string]struct{})
			for rows.Next() {
				var id string
				if err := rows.Scan(&id); err != nil {
					return errors.InternalErr.Wrap(err)
				}
				currentIDs[id] = struct{}{}
			}

			// Determine IDs to insert and delete
			newIDs := make(map[string]struct{})
			for _, id := range dst.Values {
				newIDs[id] = struct{}{}
			}

			var toInsert []string
			var toDelete []string

			for id := range newIDs {
				if _, exists := currentIDs[id]; !exists {
					toInsert = append(toInsert, id)
				}
			}
			for id := range currentIDs {
				if _, exists := newIDs[id]; !exists {
					toDelete = append(toDelete, id)
				}
			}

			// Perform insertions
			if len(toInsert) > 0 {
				for _, id := range toInsert {
					insertQuery := fmt.Sprintf("INSERT INTO %s (id) VALUES ('%s');", dst.Table, id) //nolint:gosec
					if _, err := tx.Exec(insertQuery); err != nil {
						return errors.InternalErr.Wrap(err)
					}
				}
			}

			// Perform deletions
			if len(toDelete) > 0 {
				for _, id := range toDelete {
					deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = '%s';", dst.Table, id) //nolint:gosec
					if _, err := tx.Exec(deleteQuery); err != nil {
						return errors.InternalErr.Wrap(err)
					}
				}
			}
			logger.L(ctx).Info(fmt.Sprintf("Synced table %s: %d inserted, %d deleted", dst.Table, len(toInsert), len(toDelete)))
		}
		return nil
	}(); err != nil {
		syncErr = err
		return
	}

	logger.L(ctx).Info("complete sync constants")
}
