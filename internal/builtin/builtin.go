package builtin

import "embed"

//go:embed configs/*.yaml
var FS embed.FS
