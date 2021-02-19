package web

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
)

//go:embed dist/octant
var feBundle embed.FS

// Handler create a http handler for the web content.
func Handler() (http.Handler, error) {
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		return nil, err
	}

	// This step is needed as all the assets are served under root path.
	fsys, err := fs.Sub(feBundle, "dist/octant")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(fsys)), nil
}
