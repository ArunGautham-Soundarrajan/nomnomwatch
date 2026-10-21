package assets

import "embed"

// Static holds the files under static/, served at /static/.
//
//go:embed static
var Static embed.FS
