package web

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
)

//go:generate rice embed-go

// Handler create a http handler for the web content.
func Handler() (http.Handler, error) {
	box, err := rice.FindBox("build")
	if err != nil {
		return nil, err
	}

	return http.FileServer(box.HTTPBox()), nil
}
