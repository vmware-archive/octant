package astihttp

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"io"

	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
)

// Download is a cancellable function that downloads a src into a dst using a specific *http.Client
func Download(ctx context.Context, c *http.Client, src, dst string) (err error) {
	// Create the dst file
	var f *os.File
	if f, err = os.Create(dst); err != nil {
		return errors.Wrapf(err, "astihttp: creating file %s failed", dst)
	}
	defer f.Close()

	// Download in writer
	if err = DownloadInWriter(ctx, c, src, f); err != nil {
		return errors.Wrap(err, "astihttp: downloading in writer failed")
	}
	return
}

// DownloadInWriter is a cancellable function that downloads a src into a writer using a specific *http.Client
func DownloadInWriter(ctx context.Context, c *http.Client, src string, dst io.Writer) (err error) {
	// Send request
	var resp *http.Response
	if resp, err = c.Get(src); err != nil {
		return errors.Wrapf(err, "astihttp: getting %s failed", src)
	}
	defer resp.Body.Close()

	// Validate status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("astihttp: getting %s returned %d status code", src, resp.StatusCode)
	}

	// Copy
	if _, err = astiio.Copy(ctx, resp.Body, dst); err != nil {
		return errors.Wrapf(err, "astihttp: copying content from %s to writer failed", src)
	}
	return
}
