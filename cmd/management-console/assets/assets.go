// +build !dist

package assets

import (
	"github.com/shurcooL/httpfs/filter"
	"net/http"
	"os"
	"path/filepath"
)

// Assets contains the bundled web assets
var Assets http.FileSystem

func init() {
	path := ":"
	if projectRoot := os.Getenv("PROJECT_ROOT"); projectRoot != "" {
		path = filepath.Join(projectRoot, "assets")
	}
	Assets = http.Dir(path)

	// Don't include Go files
	Assets = filter.Skip(Assets, filter.FilesWithExtensions(".go"))
}
