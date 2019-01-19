package astios

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"os"

	"github.com/pkg/errors"
)

// Checksum computes the checksum of a file
func Checksum(path string) (checksum string, err error) {
	// Open executable
	var f *os.File
	if f, err = os.Open(path); err != nil {
		err = errors.Wrapf(err, "opening %s failed", path)
		return
	}
	defer f.Close()

	// Compute checksum
	var h = sha1.New()
	if _, err = io.Copy(h, f); err != nil {
		err = errors.Wrap(err, "copying file to hasher failed")
		return
	}
	checksum = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}
