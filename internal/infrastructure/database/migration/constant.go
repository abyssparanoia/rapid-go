package migration

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"strings"

	constant_files "github.com/abyssparanoia/rapid-go/db/main/constants"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

// constantData represents the structure of a YAML file for a constants table.
type constantData struct {
	IDs []string `yaml:"ids"`
}

func syncConstantByYaml( //nolint:gocognit
	ctx context.Context,
	tx *sql.Tx,
	yamlFileEntry fs.DirEntry,
) error {
	tableName := strings.TrimSuffix(yamlFileEntry.Name(), ".yaml")
	tableName = strings.TrimSuffix(tableName, ".yml")

	fileContent, err := constant_files.EmbedConstants.ReadFile(yamlFileEntry.Name())
	if err != nil {
		return errors.InternalErr.Wrap(err)
	}

	var data constantData
	if err = yaml.Unmarshal(fileContent, &data); err != nil {
		return errors.InternalErr.Wrap(err)
	}

	// Fetch current IDs from the database
	query := fmt.Sprintf("SELECT id FROM %s", tableName) //nolint:gosec
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
	for _, id := range data.IDs {
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
		insertQuery := fmt.Sprintf("INSERT INTO %s (id) VALUES (?)", tableName) //nolint:gosec
		for _, id := range toInsert {
			if _, err := tx.Exec(insertQuery, id); err != nil {
				return errors.InternalErr.Wrap(err)
			}
		}
	}

	// Perform deletions
	if len(toDelete) > 0 {
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName) //nolint:gosec
		for _, id := range toDelete {
			if _, err := tx.Exec(deleteQuery, id); err != nil {
				return errors.InternalErr.Wrap(err)
			}
		}
	}

	logger.L(ctx).Info(fmt.Sprintf("Synced table %s: %d inserted, %d deleted", tableName, len(toInsert), len(toDelete)))

	return nil
}
