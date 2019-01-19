package astios

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
)

// Copy is a cross partitions cancellable copy
// If src is a file, dst must be the full path to file once copied
// If src is a dir, dst must be the full path to the dir once copied
func Copy(ctx context.Context, src, dst string) (err error) {
	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Stat src
	var statSrc os.FileInfo
	if statSrc, err = os.Stat(src); err != nil {
		err = errors.Wrapf(err, "stating %s failed", src)
		return
	}

	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Dir
	if statSrc.IsDir() {
		if err = filepath.Walk(src, func(path string, info os.FileInfo, errWalk error) error {
			// Check error
			if errWalk != nil {
				return errWalk
			}

			// Do not process root
			if src == path {
				return nil
			}

			// Copy
			var p = filepath.Join(dst, strings.TrimPrefix(path, filepath.Clean(src)))
			if errCopy := Copy(ctx, path, p); errCopy != nil {
				return errors.Wrapf(err, "copying %s to %s failed", path, p)
			}
			return nil
		}); err != nil {
			return
		}
		return
	}

	// Open the source file
	var srcFile *os.File
	if srcFile, err = os.Open(src); err != nil {
		err = errors.Wrapf(err, "opening %s failed", src)
		return
	}
	defer srcFile.Close()

	// Create the destination folder
	if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		err = errors.Wrapf(err, "mkdirall %s failed", filepath.Dir(dst))
		return
	}

	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Create the destination file
	var dstFile *os.File
	if dstFile, err = os.Create(dst); err != nil {
		err = errors.Wrapf(err, "creating %s failed", dst)
		return
	}
	defer dstFile.Close()

	// Chmod using os.chmod instead of file.Chmod
	if err = os.Chmod(dst, statSrc.Mode()); err != nil {
		err = errors.Wrapf(err, "chmod %s %s failed", dst, statSrc.Mode())
		return
	}

	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Copy the content
	if _, err = astiio.Copy(ctx, srcFile, dstFile); err != nil {
		err = errors.Wrapf(err, "copying content of %s to %s failed", srcFile.Name(), dstFile.Name())
		return
	}

	return
}
