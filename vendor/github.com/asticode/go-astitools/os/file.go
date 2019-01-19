package astios

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// TempFile writes a content to a temp file
func TempFile(i []byte) (name string, err error) {
	// Create temp file
	var f *os.File
	if f, err = ioutil.TempFile(os.TempDir(), "astitools"); err != nil {
		err = errors.Wrap(err, "creating temp file failed")
		return
	}
	name = f.Name()
	defer f.Close()

	// Write
	if _, err = f.Write(i); err != nil {
		err = errors.Wrapf(err, "writing to %s failed", name)
		return
	}
	return
}
