package migration_files

import "embed"

//go:embed *.sql
var EmbedMigrations embed.FS
