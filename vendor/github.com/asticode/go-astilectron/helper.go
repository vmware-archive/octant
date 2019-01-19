package astilectron

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/archive"
	"github.com/asticode/go-astitools/context"
	"github.com/asticode/go-astitools/http"
	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
)

// Download is a cancellable function that downloads a src into a dst using a specific *http.Client and cleans up on
// failed downloads
func Download(ctx context.Context, c *http.Client, src, dst string) (err error) {
	// Log
	astilog.Debugf("Downloading %s into %s", src, dst)

	// Destination already exists
	if _, err = os.Stat(dst); err == nil {
		astilog.Debugf("%s already exists, skipping download...", dst)
		return
	} else if !os.IsNotExist(err) {
		return errors.Wrapf(err, "stating %s failed", dst)
	}
	err = nil

	// Clean up on error
	defer func(err *error) {
		if *err != nil || ctx.Err() != nil {
			astilog.Debugf("Removing %s...", dst)
			os.Remove(dst)
		}
	}(&err)

	// Make sure the dst directory  exists
	if err = os.MkdirAll(filepath.Dir(dst), 0775); err != nil {
		return errors.Wrapf(err, "mkdirall %s failed", filepath.Dir(dst))
	}

	// Download
	if err = astihttp.Download(ctx, c, src, dst); err != nil {
		return errors.Wrap(err, "astihttp.Download failed")
	}
	return
}

// Disembed is a cancellable disembed of an src to a dst using a custom Disembedder
func Disembed(ctx context.Context, d Disembedder, src, dst string) (err error) {
	// Log
	astilog.Debugf("Disembedding %s into %s...", src, dst)

	// No need to disembed
	if _, err = os.Stat(dst); err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "stating %s failed", dst)
	} else if err == nil {
		astilog.Debugf("%s already exists, skipping disembed...", dst)
		return
	}
	err = nil

	// Clean up on error
	defer func(err *error) {
		if *err != nil || ctx.Err() != nil {
			astilog.Debugf("Removing %s...", dst)
			os.Remove(dst)
		}
	}(&err)

	// Make sure directory exists
	var dirPath = filepath.Dir(dst)
	astilog.Debugf("Creating %s", dirPath)
	if err = os.MkdirAll(dirPath, 0755); err != nil {
		return errors.Wrapf(err, "mkdirall %s failed", dirPath)
	}

	// Create dst
	var f *os.File
	astilog.Debugf("Creating %s", dst)
	if f, err = os.Create(dst); err != nil {
		return errors.Wrapf(err, "creating %s failed", dst)
	}
	defer f.Close()

	// Disembed
	var b []byte
	astilog.Debugf("Disembedding %s", src)
	if b, err = d(src); err != nil {
		return errors.Wrapf(err, "disembedding %s failed", src)
	}

	// Copy
	astilog.Debugf("Copying disembedded data to %s", dst)
	if _, err = astiio.Copy(ctx, bytes.NewReader(b), f); err != nil {
		return errors.Wrapf(err, "copying disembedded data into %s failed", dst)
	}
	return
}

// Unzip unzips a src into a dst.
// Possible src formats are /path/to/zip.zip or /path/to/zip.zip/internal/path.
func Unzip(ctx context.Context, src, dst string) (err error) {
	// Clean up on error
	defer func(err *error) {
		if *err != nil || ctx.Err() != nil {
			astilog.Debugf("Removing %s...", dst)
			os.RemoveAll(dst)
		}
	}(&err)

	// Unzipping
	astilog.Debugf("Unzipping %s into %s", src, dst)
	if err = astiarchive.Unzip(ctx, src, dst); err != nil {
		err = errors.Wrapf(err, "unzipping %s into %s failed", src, dst)
		return
	}
	return
}

// PtrBool transforms a bool into a *bool
func PtrBool(i bool) *bool {
	return &i
}

// PtrFloat transforms a float64 into a *float64
func PtrFloat(i float64) *float64 {
	return &i
}

// PtrInt transforms an int into an *int
func PtrInt(i int) *int {
	return &i
}

// PtrInt64 transforms an int64 into an *int64
func PtrInt64(i int64) *int64 {
	return &i
}

// PtrStr transforms a string into a *string
func PtrStr(i string) *string {
	return &i
}

// synchronousFunc executes a function, blocks until it has received a specific event or the canceller has been
// cancelled and returns the corresponding event
func synchronousFunc(c *asticontext.Canceller, l listenable, fn func(), eventNameDone string) (e Event) {
	var ctx, cancel = c.NewContext()
	defer cancel()
	l.On(eventNameDone, func(i Event) (deleteListener bool) {
		e = i
		cancel()
		return true
	})
	fn()
	<-ctx.Done()
	return
}

// synchronousEvent sends an event, blocks until it has received a specific event or the canceller has been cancelled
// and returns the corresponding event
func synchronousEvent(c *asticontext.Canceller, l listenable, w *writer, i Event, eventNameDone string) (o Event, err error) {
	o = synchronousFunc(c, l, func() {
		if err = w.write(i); err != nil {
			err = errors.Wrapf(err, "writing %+v event failed", i)
			return
		}
		return
	}, eventNameDone)
	return
}
