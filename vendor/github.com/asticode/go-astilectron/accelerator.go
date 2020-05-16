package astilectron

import (
	"strings"
)

// Accelerator separator
const acceleratorSeparator = "+"

// Accelerator represents an accelerator
// https://github.com/electron/electron/blob/v1.8.1/docs/api/accelerator.md
type Accelerator []string

// NewAccelerator creates a new accelerator
func NewAccelerator(items ...string) (a *Accelerator) {
	a = &Accelerator{}
	for _, i := range items {
		*a = append(*a, i)
	}
	return
}

// MarshalText implements the encoding.TextMarshaler interface
func (a *Accelerator) MarshalText() ([]byte, error) {
	return []byte(strings.Join(*a, acceleratorSeparator)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (a *Accelerator) UnmarshalText(b []byte) error {
	*a = strings.Split(string(b), acceleratorSeparator)
	return nil
}
