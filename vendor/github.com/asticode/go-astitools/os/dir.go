package astios

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// TempDir creates a temp dir
func TempDir(prefix string) (path string, err error) {
	// Create temp file
	var f *os.File
	if f, err = ioutil.TempFile(os.TempDir(), prefix); err != nil {
		err = errors.Wrap(err, "creating temporary file failed")
		return
	}
	path = f.Name()

	// Close temp file
	if err = f.Close(); err != nil {
		err = errors.Wrapf(err, "closing file %s failed", path)
		return
	}

	// Delete temp file
	if err = os.Remove(path); err != nil {
		err = errors.Wrapf(err, "removing %s failed", path)
		return
	}

	// Create temp dir
	if err = os.MkdirAll(path, 0755); err != nil {
		err = errors.Wrapf(err, "mkdirall of %s failed", path)
		return
	}
	return
}
