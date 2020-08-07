/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package tsgen

import (
	"bytes"
	"os/exec"
)

// TSFormatter formats typescript.
type TSFormatter struct{}

// NewTSFormatter creates an instance of the typescript formatter.
func NewTSFormatter() *TSFormatter {
	tf := &TSFormatter{}
	return tf
}

// Format formats the supplied typescript.
func (tf *TSFormatter) Format(in []byte) ([]byte, error) {
	cmd := exec.Command("prettier", "--stdin-filepath", "file.ts", "--single-quote")

	r := bytes.NewReader(in)
	cmd.Stdin = r

	return cmd.Output()
}
