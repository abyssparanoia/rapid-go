package constant_files

import "embed"

//go:embed *.yaml
var EmbedConstants embed.FS
