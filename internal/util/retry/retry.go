/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package retry

import (
	"math/rand"
	"time"
)

// Retry retries func `f` `attempts` times with `sleep` between retries.
func Retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(StopError); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

// StopError is an error that will cause the retry loop to stop
type StopError struct {
	error
}

// Stop stops the retry with an error.
func Stop(err error) *StopError {
	return &StopError{error: err}
}
