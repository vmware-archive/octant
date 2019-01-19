package astiarchive

import (
	"archive/zip"
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
)

// Zip zips a src into a dst
// dstRoot can be used to create root directories so that the content is zipped in /path/to/zip.zip/root/path
func Zip(ctx context.Context, src, dst, dstRoot string) (err error) {
	// Create destination file
	var dstFile *os.File
	if dstFile, err = os.Create(dst); err != nil {
		return err
	}
	defer dstFile.Close()

	// Create zip writer
	var zw = zip.NewWriter(dstFile)
	defer zw.Close()

	// Walk
	filepath.Walk(src, func(path string, info os.FileInfo, e1 error) (e2 error) {
		// Process error
		if e1 != nil {
			return e1
		}

		// Init header
		var h *zip.FileHeader
		if h, e2 = zip.FileInfoHeader(info); e2 != nil {
			return
		}

		// Set header info
		h.Name = filepath.Join(dstRoot, strings.TrimPrefix(path, src))
		if info.IsDir() {
			h.Name += "/"
		} else {
			h.Method = zip.Deflate
		}

		// Create writer
		var w io.Writer
		if w, e2 = zw.CreateHeader(h); e2 != nil {
			return
		}

		// If path is dir, stop here
		if info.IsDir() {
			return
		}

		// Open path
		var walkFile *os.File
		if walkFile, e2 = os.Open(path); e2 != nil {
			return
		}
		defer walkFile.Close()

		// Copy
		if _, e2 = astiio.Copy(ctx, walkFile, w); e2 != nil {
			return
		}
		return
	})
	return
}

// Unzip unzips a src into a dst
// Possible src formats are /path/to/zip.zip or /path/to/zip.zip/internal/path if you only want to unzip files in
// /internal/path in the .zip archive
func Unzip(ctx context.Context, src, dst string) (err error) {
	// Parse src path
	var split = strings.Split(src, ".zip")
	src = split[0] + ".zip"
	var internalPath string
	if len(split) >= 2 {
		internalPath = split[1]
	}

	// Open overall reader
	var r *zip.ReadCloser
	if r, err = zip.OpenReader(src); err != nil {
		return errors.Wrapf(err, "astiarchive: opening overall zip reader on %s failed", src)
	}
	defer r.Close()

	// Loop through files to determine their type
	var dirs, files, symlinks = make(map[string]*zip.File), make(map[string]*zip.File), make(map[string]*zip.File)
	for _, f := range r.File {
		// Validate internal path
		var n = string(os.PathSeparator) + f.Name
		if internalPath != "" && !strings.HasPrefix(n, internalPath) {
			continue
		}
		var p = filepath.Join(dst, strings.TrimPrefix(n, internalPath))

		// Check file type
		if f.FileInfo().Mode()&os.ModeSymlink != 0 {
			symlinks[p] = f
		} else if f.FileInfo().IsDir() {
			dirs[p] = f
		} else {
			files[p] = f
		}
	}

	// Create dirs
	for p, f := range dirs {
		if err = os.MkdirAll(p, f.FileInfo().Mode().Perm()); err != nil {
			return errors.Wrapf(err, "astiarchive: mkdirall %s failed", p)
		}
	}

	// Create files
	for p, f := range files {
		if err = createZipFile(ctx, f, p); err != nil {
			return errors.Wrapf(err, "astiarchive: creating zip file into %s failed", p)
		}
	}

	// Create symlinks
	for p, f := range symlinks {
		if err = createZipSymlink(f, p); err != nil {
			return errors.Wrapf(err, "astiarchive: creating zip symlink into %s failed", p)
		}
	}
	return
}

func createZipFile(ctx context.Context, f *zip.File, p string) (err error) {
	// Open file reader
	var fr io.ReadCloser
	if fr, err = f.Open(); err != nil {
		return errors.Wrapf(err, "astiarchive: opening zip reader on file %s failed", f.Name)
	}
	defer fr.Close()

	// Since dirs don't always come up we make sure the directory of the file exists with default
	// file mode
	if err = os.MkdirAll(filepath.Dir(p), DefaultFileMode); err != nil {
		return errors.Wrapf(err, "astiarchive: mkdirall %s failed", filepath.Dir(p))
	}

	// Open the file
	var fl *os.File
	if fl, err = os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode().Perm()); err != nil {
		return errors.Wrapf(err, "astiarchive: opening file %s failed", p)
	}
	defer fl.Close()

	// Copy
	if _, err = astiio.Copy(ctx, fr, fl); err != nil {
		return errors.Wrapf(err, "astiarchive: copying %s into %s failed", f.Name, p)
	}
	return
}

func createZipSymlink(f *zip.File, p string) (err error) {
	// Open file reader
	var fr io.ReadCloser
	if fr, err = f.Open(); err != nil {
		return errors.Wrapf(err, "astiarchive: opening zip reader on file %s failed", f.Name)
	}
	defer fr.Close()

	// If file is a symlink we retrieve the target path that is in the content of the file
	var b []byte
	if b, err = ioutil.ReadAll(fr); err != nil {
		return errors.Wrapf(err, "astiarchive: ioutil.Readall on %s failed", f.Name)
	}

	// Create the symlink
	if err = os.Symlink(string(b), p); err != nil {
		return errors.Wrapf(err, "astiarchive: creating symlink from %s to %s failed", string(b), p)
	}
	return
}
