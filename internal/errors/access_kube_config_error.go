package errors

import (
	"fmt"
	"hash/fnv"
	"time"
)

// The main difference of this Error Struct is that for the ID is
// using non-cryptographic hash to make sure ID is uniquely identify
// go provides many algorithms https://golang.org/pkg/hash but fnv was chosen,
// because of the information here
// https://softwareengineering.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed
// fvn provides a good compromise between performance and number of collisions

const OctantWrongCertificateError = "OctantWrongCertificateError"

type WrongCertificateError struct {
	id        string    `json:"id"`
	timestamp time.Time `json:"timestamp"`
	err       error     `json:"error"`
}

func NewWrongCertificateError(err error) *WrongCertificateError {
	data := fnv.New64().Sum([]byte(OctantWrongCertificateError + ": " + err.Error()))

	return &WrongCertificateError{
		err:       err,
		timestamp: time.Now(),
		id:        fmt.Sprintf("%x", data),
	}
}

var _ InternalError = (*WrongCertificateError)(nil)

func (o *WrongCertificateError) Name() string {
	return OctantWrongCertificateError
}

// ID returns the error unique ID.
func (o *WrongCertificateError) ID() string {
	return o.id
}

// Timestamp returns the error timestamp.
func (o *WrongCertificateError) Timestamp() time.Time {
	return o.timestamp
}

// Error returns an error string.
func (o *WrongCertificateError) Error() string {
	return fmt.Sprintf("%s", o.err)
}
