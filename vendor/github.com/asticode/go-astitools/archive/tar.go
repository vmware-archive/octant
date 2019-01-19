package astiarchive

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"os"

	"io"

	"path/filepath"

	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
)

// Untar untars a src into a dst
func Untar(ctx context.Context, src, dst string) (err error) {
	// Open src
	var srcFile *os.File
	if srcFile, err = os.Open(src); err != nil {
		return errors.Wrapf(err, "astiarchive: opening %s failed", src)
	}
	defer srcFile.Close()

	// Create gzip reader
	var gzr *gzip.Reader
	if gzr, err = gzip.NewReader(srcFile); err != nil {
		return errors.Wrap(err, "astiarchive: creating gzip reader failed")
	}
	defer gzr.Close()

	// Loop through tar entries
	tr := tar.NewReader(gzr)
	for {
		// Get next entry
		var h *tar.Header
		if h, err = tr.Next(); err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return errors.Wrap(err, "astiarchive: getting next header failed")
		}

		// No header
		if h == nil {
			continue
		}

		// Build path
		p := filepath.Join(dst, h.Name)

		// Switch on file type
		switch h.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(p, h.FileInfo().Mode().Perm()); err != nil {
				err = errors.Wrapf(err, "astiarchive: mkdirall %s failed", p)
				return
			}
		case tar.TypeReg:
			if err = createTarFile(ctx, p, h, tr); err != nil {
				err = errors.Wrapf(err, "astiarchive: creating tar file into %s failed", p)
				return
			}
		}
	}
	return
}

func createTarFile(ctx context.Context, p string, h *tar.Header, tr *tar.Reader) (err error) {
	// Sometimes the dir that will contain the file has not yet been processed in the tar ball, therefore we need to create it
	if err = os.MkdirAll(filepath.Dir(p), DefaultFileMode); err != nil {
		err = errors.Wrapf(err, "astiarchive: mkdirall %s failed", filepath.Dir(p))
		return
	}

	// Open file
	var f *os.File
	if f, err = os.OpenFile(p, os.O_TRUNC|os.O_CREATE|os.O_RDWR, h.FileInfo().Mode().Perm()); err != nil {
		err = errors.Wrap(err, "astiarchive: opening file failed")
		return
	}
	defer f.Close()

	// Copy
	if _, err = astiio.Copy(ctx, tr, f); err != nil {
		err = errors.Wrap(err, "astiarchive: copying content failed")
		return
	}
	return
}
