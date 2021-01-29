package web

import (
	"mime"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
)

//go:generate rice embed-go

// Handler create a http handler for the web content.
func Handler() (http.Handler, error) {
	box, err := rice.FindBox("dist/octant")
	if err != nil {
		return nil, err
	}

	err = mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		return nil, err
	}

	return http.FileServer(box.HTTPBox()), nil
}
