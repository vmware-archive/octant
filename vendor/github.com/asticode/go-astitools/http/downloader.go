package astihttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
)

// Downloader represents a downloader
type Downloader struct {
	bp              *sync.Pool
	busyWorkers     int
	cond            *sync.Cond
	ignoreErrors    bool
	mc              *sync.Mutex // Locks cond
	mw              *sync.Mutex // Locks busyWorkers
	numberOfWorkers int
	s               *Sender
}

// DownloaderFunc represents a downloader func
// It's its responsibility to close the reader
type DownloaderFunc func(ctx context.Context, idx int, src string, r io.ReadCloser) error

// DownloaderOptions represents downloader options
type DownloaderOptions struct {
	IgnoreErrors    bool
	NumberOfWorkers int
	Sender          SenderOptions
}

// NewDownloader creates a new downloader
func NewDownloader(o DownloaderOptions) (d *Downloader) {
	d = &Downloader{
		bp:              &sync.Pool{New: func() interface{} { return &bytes.Buffer{} }},
		ignoreErrors:    o.IgnoreErrors,
		mc:              &sync.Mutex{},
		mw:              &sync.Mutex{},
		numberOfWorkers: o.NumberOfWorkers,
		s:               NewSender(o.Sender),
	}
	d.cond = sync.NewCond(d.mc)
	if d.numberOfWorkers == 0 {
		d.numberOfWorkers = 1
	}
	return
}

// Download downloads in parallel a set of src paths and executes a custom callback on each downloaded buffers
func (d *Downloader) Download(parentCtx context.Context, paths []string, fn DownloaderFunc) (err error) {
	// Init
	ctx, cancel := context.WithCancel(parentCtx)
	m := &sync.Mutex{} // Locks err
	gwg := &sync.WaitGroup{}
	gwg.Add(len(paths))
	lwg := &sync.WaitGroup{}

	// Loop through src paths
	var idx int
	for idx < len(paths) {
		// Check context
		if parentCtx.Err() != nil {
			m.Lock()
			err = errors.Wrap(err, "astihttp: context error")
			m.Unlock()
			return
		}

		// Lock cond here in case a worker finishes between checking the number of busy workers and the if statement
		d.cond.L.Lock()

		// Check if a worker is available
		var ok bool
		d.mw.Lock()
		if ok = d.numberOfWorkers > d.busyWorkers; ok {
			d.busyWorkers++
		}
		d.mw.Unlock()

		// No worker is available
		if !ok {
			d.cond.Wait()
			d.cond.L.Unlock()
			continue
		}
		d.cond.L.Unlock()

		// Check error
		m.Lock()
		if err != nil {
			m.Unlock()
			lwg.Wait()
			return
		}
		m.Unlock()

		// Download
		go func(idx int) {
			lwg.Add(1)
			if errR := d.download(ctx, idx, paths[idx], fn, gwg); errR != nil {
				m.Lock()
				if err == nil {
					err = errR
				}
				m.Unlock()
				cancel()
			}
			lwg.Done()
		}(idx)
		idx++
	}
	gwg.Wait()
	return
}

type readCloser struct {
	b  *bytes.Buffer
	bp *sync.Pool
}

func newReadCloser(b *bytes.Buffer, bp *sync.Pool) *readCloser {
	return &readCloser{
		b:  b,
		bp: bp,
	}
}

// Read implements the io.Reader interface
func (c readCloser) Read(p []byte) (n int, err error) {
	return c.b.Read(p)
}

// Close implements the io.Closer interface
func (c readCloser) Close() error {
	c.b.Reset()
	c.bp.Put(c.b)
	return nil
}

func (d *Downloader) download(ctx context.Context, idx int, path string, fn DownloaderFunc, wg *sync.WaitGroup) (err error) {
	// Update wait group and worker status
	defer func() {
		// Update worker status
		d.mw.Lock()
		d.busyWorkers--
		d.mw.Unlock()

		// Broadcast
		d.cond.L.Lock()
		d.cond.Broadcast()
		d.cond.L.Unlock()

		// Update wait group
		wg.Done()
	}()

	// Create request
	var r *http.Request
	if r, err = http.NewRequest(http.MethodGet, path, nil); err != nil {
		return errors.Wrapf(err, "astihttp: creating GET request to %s failed", path)
	}

	// Send request
	var resp *http.Response
	if resp, err = d.s.Send(r); err != nil {
		return errors.Wrapf(err, "astihttp: sending GET request to %s failed", path)
	}
	defer resp.Body.Close()

	// Validate status code
	buf := newReadCloser(d.bp.Get().(*bytes.Buffer), d.bp)
	if resp.StatusCode != http.StatusOK {
		errS := fmt.Errorf("astihttp: sending GET request to %s returned %d status code", path, resp.StatusCode)
		if !d.ignoreErrors {
			return errS
		} else {
			astilog.Error(errors.Wrap(errS, "astihttp: ignoring error"))
		}
	} else {
		// Copy body
		if _, err = astiio.Copy(ctx, resp.Body, buf.b); err != nil {
			return errors.Wrap(err, "astihttp: copying resp.Body to buf.b failed")
		}
	}

	// Custom callback
	if err = fn(ctx, idx, path, buf); err != nil {
		return errors.Wrapf(err, "astihttp: custom callback on %s failed", path)
	}
	return
}

// DownloadInDirectory downloads in parallel a set of src paths and saves them in a dst directory
func (d *Downloader) DownloadInDirectory(ctx context.Context, dst string, paths ...string) error {
	return d.Download(ctx, paths, func(ctx context.Context, idx int, path string, r io.ReadCloser) (err error) {
		// Make sure to close the reader
		defer r.Close()

		// Make sure destination directory exists
		if err = os.MkdirAll(dst, 0700); err != nil {
			err = errors.Wrapf(err, "astihttp: mkdirall %s failed", dst)
			return
		}

		// Create destination file
		var f *os.File
		dst := filepath.Join(dst, filepath.Base(path))
		if f, err = os.Create(dst); err != nil {
			err = errors.Wrapf(err, "astihttp: creating %s failed", dst)
			return
		}
		defer f.Close()

		// Copy
		if _, err = astiio.Copy(ctx, r, f); err != nil {
			err = errors.Wrapf(err, "astihttp: copying content to %s failed", dst)
			return
		}
		return
	})
}

type chunk struct {
	idx  int
	r    io.ReadCloser
	path string
}

// DownloadInWriter downloads in parallel a set of src paths and concatenates them in order in a writer
func (d *Downloader) DownloadInWriter(ctx context.Context, w io.Writer, paths ...string) (err error) {
	// Download
	var cs []chunk
	var m sync.Mutex // Locks cs
	var requiredIdx int
	err = d.Download(ctx, paths, func(ctx context.Context, idx int, path string, r io.ReadCloser) (err error) {
		// Lock
		m.Lock()
		defer m.Unlock()

		// Check where to insert chunk
		var idxInsert = -1
		for idxChunk := 0; idxChunk < len(cs); idxChunk++ {
			if idx < cs[idxChunk].idx {
				idxInsert = idxChunk
				break
			}
		}

		// Create chunk
		c := chunk{
			idx:  idx,
			path: path,
			r:    r,
		}

		// Add chunk
		if idxInsert > -1 {
			cs = append(cs[:idxInsert], append([]chunk{c}, cs[idxInsert:]...)...)
		} else {
			cs = append(cs, c)
		}

		// Loop through chunks
		for idxChunk := 0; idxChunk < len(cs); idxChunk++ {
			// Get chunk
			c := cs[idxChunk]

			// The chunk should be copied
			if c.idx == requiredIdx {
				// Copy chunk content
				_, err = astiio.Copy(ctx, c.r, w)

				// Make sure the reader is closed
				c.r.Close()

				// Remove chunk
				requiredIdx++
				cs = append(cs[:idxChunk], cs[idxChunk+1:]...)
				idxChunk--

				// Check error now so that chunk is still removed and reader is closed
				if err != nil {
					err = errors.Wrapf(err, "astihttp: copying chunk #%d to dst failed", c.idx)
					return
				}
			}
		}
		return
	})

	// Make sure to close all readers
	for _, c := range cs {
		c.r.Close()
	}
	return
}

// DownloadInFile downloads in parallel a set of src paths and concatenates them in order in a writer
func (d *Downloader) DownloadInFile(ctx context.Context, dst string, paths ...string) (err error) {
	// Make sure destination directory exists
	if err = os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
		err = errors.Wrapf(err, "astihttp: mkdirall %s failed", filepath.Dir(dst))
		return
	}

	// Create destination file
	var f *os.File
	if f, err = os.Create(dst); err != nil {
		err = errors.Wrapf(err, "astihttp: creating %s failed", dst)
		return
	}
	defer f.Close()

	// Download in writer
	return d.DownloadInWriter(ctx, f, paths...)
}
