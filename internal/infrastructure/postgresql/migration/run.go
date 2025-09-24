package migration

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	constant_files "github.com/abyssparanoia/rapid-go/db/postgresql/constants"
	migration_files "github.com/abyssparanoia/rapid-go/db/postgresql/migrations"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/postgresql"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger/logger_field"
	"github.com/caarlos0/env/v11"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gopkg.in/yaml.v3"
)

func RunNewFile(fileName string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := goose.Create(nil, filepath.Join(dir, "internal/infrastructure/postgresql/migration/files"), fileName, "sql"); err != nil {
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

	databaseCli := postgresql.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, true)

	goose.SetBaseFS(migration_files.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
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

	databaseCli := postgresql.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, true)

	ctx := context.Background()
	tables, err := listTables(ctx, databaseCli.DB, "public")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("./db/postgresql/schema.sql")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, table := range tables {
		ddl, err := buildTableDefinition(ctx, databaseCli.DB, "public", table)
		if err != nil {
			panic(err)
		}
		if _, err := file.WriteString(ddl); err != nil {
			panic(err)
		}
	}

	logger.Info("complete extracting database schema")
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

	databaseCli := postgresql.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, true)
	runSyncConstantsWithContext(ctx, databaseCli)
}

func runSyncConstantsWithContext( //nolint:gocognit
	ctx context.Context,
	databaseCli *postgresql.Client,
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

func listTables(ctx context.Context, db *sql.DB, schema string) ([]string, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT tablename
		FROM pg_catalog.pg_tables
		WHERE schemaname = $1
		ORDER BY tablename`,
		schema,
	)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, errors.InternalErr.Wrap(err)
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}

	sort.Strings(tables)
	return tables, nil
}

func buildTableDefinition(ctx context.Context, db *sql.DB, schema, table string) (string, error) {
	columns, err := fetchColumns(ctx, db, schema, table)
	if err != nil {
		return "", err
	}

	constraints, err := fetchConstraints(ctx, db, schema, table)
	if err != nil {
		return "", err
	}

	indexes, err := fetchIndexes(ctx, db, schema, table)
	if err != nil {
		return "", err
	}

	var lines []string
	for _, column := range columns {
		lines = append(lines, fmt.Sprintf("    %s", column))
	}
	for _, constraint := range constraints {
		lines = append(lines, fmt.Sprintf("    CONSTRAINT %s %s", pq.QuoteIdentifier(constraint.Name), constraint.Definition))
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("CREATE TABLE %s.%s (\n", pq.QuoteIdentifier(schema), pq.QuoteIdentifier(table)))
	buf.WriteString(strings.Join(lines, ",\n"))
	buf.WriteString("\n);\n\n")

	for _, index := range indexes {
		buf.WriteString(index)
		buf.WriteString(";\n\n")
	}

	return buf.String(), nil
}

type constraintInfo struct {
	Name       string
	Definition string
}

func fetchColumns(ctx context.Context, db *sql.DB, schema, table string) ([]string, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT
			a.attname,
			pg_catalog.format_type(a.atttypid, a.atttypmod) AS data_type,
			pg_catalog.pg_get_expr(d.adbin, d.adrelid, true) AS column_default,
			a.attnotnull,
			a.attidentity
		FROM pg_catalog.pg_attribute a
		LEFT JOIN pg_catalog.pg_attrdef d ON a.attrelid = d.adrelid AND a.attnum = d.adnum
		WHERE a.attrelid = $1::regclass
		  AND a.attnum > 0
		  AND NOT a.attisdropped
		ORDER BY a.attnum`,
		fmt.Sprintf("%s.%s", pq.QuoteIdentifier(schema), pq.QuoteIdentifier(table)),
	)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var (
			name       string
			dataType   string
			defaultVal sql.NullString
			notNull    bool
			identity   sql.NullString
		)
		if err := rows.Scan(&name, &dataType, &defaultVal, &notNull, &identity); err != nil {
			return nil, errors.InternalErr.Wrap(err)
		}

		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("%s %s", pq.QuoteIdentifier(name), dataType))

		if identity.Valid && identity.String != "" {
			if identity.String == "a" {
				builder.WriteString(" GENERATED ALWAYS AS IDENTITY")
			} else {
				builder.WriteString(" GENERATED BY DEFAULT AS IDENTITY")
			}
		} else if defaultVal.Valid {
			builder.WriteString(fmt.Sprintf(" DEFAULT %s", defaultVal.String))
		}

		if notNull {
			builder.WriteString(" NOT NULL")
		}

		columns = append(columns, builder.String())
	}

	if err := rows.Err(); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}

	return columns, nil
}

func fetchConstraints(ctx context.Context, db *sql.DB, schema, table string) ([]constraintInfo, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT c.conname, c.contype, pg_catalog.pg_get_constraintdef(c.oid, true) AS definition
		FROM pg_catalog.pg_constraint c
		WHERE c.conrelid = $1::regclass
		  AND c.contype IN ('p', 'u', 'f', 'c')
		ORDER BY c.contype, c.conname`,
		fmt.Sprintf("%s.%s", pq.QuoteIdentifier(schema), pq.QuoteIdentifier(table)),
	)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	defer rows.Close()

	var constraints []constraintInfo
	for rows.Next() {
		var info constraintInfo
		var constraintType string
		if err := rows.Scan(&info.Name, &constraintType, &info.Definition); err != nil {
			return nil, errors.InternalErr.Wrap(err)
		}
		constraints = append(constraints, info)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}

	return constraints, nil
}

func fetchIndexes(ctx context.Context, db *sql.DB, schema, table string) ([]string, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT pg_catalog.pg_get_indexdef(i.oid, 0, true) AS index_definition
		FROM pg_catalog.pg_class t
		INNER JOIN pg_catalog.pg_namespace n ON n.oid = t.relnamespace
		INNER JOIN pg_catalog.pg_index ix ON t.oid = ix.indrelid
		INNER JOIN pg_catalog.pg_class i ON i.oid = ix.indexrelid
		WHERE n.nspname = $1
		  AND t.relname = $2
		  AND NOT ix.indisprimary
		  AND NOT EXISTS (
			  SELECT 1
			  FROM pg_catalog.pg_constraint c
			  WHERE c.conindid = i.oid
			)
		ORDER BY i.relname`,
		schema,
		table,
	)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	defer rows.Close()

	var indexes []string
	for rows.Next() {
		var definition string
		if err := rows.Scan(&definition); err != nil {
			return nil, errors.InternalErr.Wrap(err)
		}
		indexes = append(indexes, definition)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}

	return indexes, nil
}
