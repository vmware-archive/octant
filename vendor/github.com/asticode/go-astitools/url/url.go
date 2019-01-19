package astiurl

import (
	"net/url"
	"path/filepath"

	"github.com/pkg/errors"
)

// Parse parses an URL (files included)
func Parse(i string) (o *url.URL, err error) {
	// Basic parse
	if o, err = url.Parse(i); err != nil {
		err = errors.Wrapf(err, "basic parsing of url %s failed", i)
		return
	}

	// File
	if o.Scheme == "" {
		// Get absolute path
		if i, err = filepath.Abs(i); err != nil {
			err = errors.Wrapf(err, "getting absolute path of %s failed", i)
			return
		}

		// Set url
		o = &url.URL{Path: i, Scheme: "file"}
	}
	return
}
