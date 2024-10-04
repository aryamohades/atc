package static

import (
	"embed"

	"github.com/benbjohnson/hashfs"
)

//go:embed dist html
var fs embed.FS

var FS = hashfs.NewFS(fs)
