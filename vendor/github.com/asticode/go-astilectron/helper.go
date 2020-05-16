package astilectron

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/asticode/go-astikit"
)

// Download is a cancellable function that downloads a src into a dst using a specific *http.Client and cleans up on
// failed downloads
func Download(ctx context.Context, l astikit.SeverityLogger, d *astikit.HTTPDownloader, src, dst string) (err error) {
	// Log
	l.Debugf("Downloading %s into %s", src, dst)

	// Destination already exists
	if _, err = os.Stat(dst); err == nil {
		l.Debugf("%s already exists, skipping download...", dst)
		return
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stating %s failed: %w", dst, err)
	}
	err = nil

	// Clean up on error
	defer func(err *error) {
		if *err != nil || ctx.Err() != nil {
			l.Debugf("Removing %s...", dst)
			os.Remove(dst)
		}
	}(&err)

	// Make sure the dst directory  exists
	if err = os.MkdirAll(filepath.Dir(dst), 0775); err != nil {
		return fmt.Errorf("mkdirall %s failed: %w", filepath.Dir(dst), err)
	}

	// Download
	if err = d.DownloadInFile(ctx, dst, astikit.HTTPDownloaderSrc{URL: src}); err != nil {
		return fmt.Errorf("DownloadInFile failed: %w", err)
	}
	return
}

// Disembed is a cancellable disembed of an src to a dst using a custom Disembedder
func Disembed(ctx context.Context, l astikit.SeverityLogger, d Disembedder, src, dst string) (err error) {
	// Log
	l.Debugf("Disembedding %s into %s...", src, dst)

	// No need to disembed
	if _, err = os.Stat(dst); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("stating %s failed: %w", dst, err)
	} else if err == nil {
		l.Debugf("%s already exists, skipping disembed...", dst)
		return
	}
	err = nil

	// Clean up on error
	defer func(err *error) {
		if *err != nil || ctx.Err() != nil {
			l.Debugf("Removing %s...", dst)
			os.Remove(dst)
		}
	}(&err)

	// Make sure directory exists
	var dirPath = filepath.Dir(dst)
	l.Debugf("Creating %s", dirPath)
	if err = os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("mkdirall %s failed: %w", dirPath, err)
	}

	// Create dst
	var f *os.File
	l.Debugf("Creating %s", dst)
	if f, err = os.Create(dst); err != nil {
		return fmt.Errorf("creating %s failed: %w", dst, err)
	}
	defer f.Close()

	// Disembed
	var b []byte
	l.Debugf("Disembedding %s", src)
	if b, err = d(src); err != nil {
		return fmt.Errorf("disembedding %s failed: %w", src, err)
	}

	// Copy
	l.Debugf("Copying disembedded data to %s", dst)
	if _, err = astikit.Copy(ctx, f, bytes.NewReader(b)); err != nil {
		return fmt.Errorf("copying disembedded data into %s failed: %w", dst, err)
	}
	return
}

// Unzip unzips a src into a dst.
// Possible src formats are /path/to/zip.zip or /path/to/zip.zip/internal/path.
func Unzip(ctx context.Context, l astikit.SeverityLogger, src, dst string) (err error) {
	// Clean up on error
	defer func(err *error) {
		if *err != nil || ctx.Err() != nil {
			l.Debugf("Removing %s...", dst)
			os.RemoveAll(dst)
		}
	}(&err)

	// Unzipping
	l.Debugf("Unzipping %s into %s", src, dst)
	if err = astikit.Unzip(ctx, dst, src); err != nil {
		err = fmt.Errorf("unzipping %s into %s failed: %w", src, dst, err)
		return
	}
	return
}

// synchronousFunc executes a function, blocks until it has received a specific event or the context has been
// cancelled and returns the corresponding event
func synchronousFunc(parentCtx context.Context, l listenable, fn func() error, eventNameDone string) (e Event, err error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	l.On(eventNameDone, func(i Event) (deleteListener bool) {
		if ctx.Err() == nil {
			e = i
		}
		cancel()
		return true
	})
	if fn != nil {
		if err = fn(); err != nil {
			return
		}
	}
	<-ctx.Done()
	return
}

// synchronousEvent sends an event, blocks until it has received a specific event or the context has been cancelled
// and returns the corresponding event
func synchronousEvent(ctx context.Context, l listenable, w *writer, i Event, eventNameDone string) (Event, error) {
	return synchronousFunc(ctx, l, func() (err error) {
		if err = w.write(i); err != nil {
			err = fmt.Errorf("writing %+v event failed: %w", i, err)
			return
		}
		return
	}, eventNameDone)
}
