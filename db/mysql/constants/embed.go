package constant_files

import _ "embed"

//go:embed constants.yaml
var EmbedConstants []byte
