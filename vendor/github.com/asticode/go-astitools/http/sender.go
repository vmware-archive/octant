package astihttp

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Sender represents an object capable of sending http requests
type Sender struct {
	client     *http.Client
	retryFunc  RetryFunc
	retryMax   int
	retrySleep time.Duration
}

// RetryFunc is a function that decides whether to retry the request
type RetryFunc func(name string, resp *http.Response) bool

func defaultRetryFunc(name string, resp *http.Response) bool {
	if resp.StatusCode >= http.StatusInternalServerError {
		astilog.Debugf("astihttp: invalid status code %d when sending %s", resp.StatusCode, name)
		return true
	}
	return false
}

// SenderOptions represents sender options
type SenderOptions struct {
	Client     *http.Client
	RetryFunc  RetryFunc
	RetryMax   int
	RetrySleep time.Duration
}

// NewSender creates a new sender
func NewSender(o SenderOptions) (s *Sender) {
	s = &Sender{
		client:     o.Client,
		retryFunc:  o.RetryFunc,
		retryMax:   o.RetryMax,
		retrySleep: o.RetrySleep,
	}
	if s.client == nil {
		s.client = &http.Client{}
	}
	if s.retryFunc == nil {
		s.retryFunc = defaultRetryFunc
	}
	return
}

// Send sends a new *http.Request
func (s *Sender) Send(req *http.Request) (resp *http.Response, err error) {
	return s.ExecWithRetry(fmt.Sprintf("%s request to %s", req.Method, req.URL), func() (*http.Response, error) { return s.client.Do(req) })
}

// ExecWithRetry handles retrying when fetching a response
// name is used for logging purposes only
func (s *Sender) ExecWithRetry(name string, fn func() (*http.Response, error)) (resp *http.Response, err error) {
	// Loop
	// We start at retryMax + 1 so that it runs at least once even if retryMax == 0
	tries := 0
	for retriesLeft := s.retryMax + 1; retriesLeft > 0; retriesLeft-- {
		// Get request name
		nr := fmt.Sprintf("%s (%d/%d)", name, s.retryMax-retriesLeft+2, s.retryMax+1)
		tries++

		// Send request
		var retry bool
		astilog.Debugf("astihttp: sending %s", nr)
		if resp, err = fn(); err != nil {
			// If error is temporary, retry
			if netError, ok := err.(net.Error); ok && netError.Temporary() {
				astilog.Debugf("astihttp: temporary error when sending %s", nr)
				retry = true
			} else {
				err = errors.Wrapf(err, "astihttp: sending %s failed", nr)
				return
			}
		}

		// Retry
		if retry || s.retryFunc(nr, resp) {
			if retriesLeft > 1 {
				astilog.Debugf("astihttp: sleeping %s and retrying... (%d retries left)", s.retrySleep, retriesLeft-1)
				time.Sleep(s.retrySleep)
			}
			continue
		}

		// Return if conditions for retrying were not met
		return
	}

	// Max retries limit reached
	err = fmt.Errorf("astihttp: sending %s failed after %d tries", name, tries)
	return
}
