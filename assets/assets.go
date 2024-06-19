package assets

import (
	"embed"
	"io/fs"

	"github.com/kettek/go-multipath/v2"
)

//go:embed **/*.png
var embedFS embed.FS
var FS multipath.FS

func init() {
	//FS.InsertFS(os.DirFS("assets"), multipath.FirstPriority) // Can be re-enabled later for user-customized assets
	sub, err := fs.Sub(embedFS, ".")
	if err != nil {
		panic(err)
	}
	FS.InsertFS(sub, multipath.LastPriority)

	FS.Walk(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return nil
	})
}
