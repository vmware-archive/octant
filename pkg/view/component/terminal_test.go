/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminal_Marshal(t *testing.T) {
	details := TerminalDetails{
		Container: "container-id",
		Command:   "/bin/bash",
		Active:    false,
	}
	input := NewTerminal("default", "term-test", "pod-name", []string{"container-id", "sidecar-container"}, details)
	actual, err := json.Marshal(input)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(filepath.Join("testdata", "terminal.json"))
	assert.NoError(t, err)

	assert.JSONEq(t, string(expected), string(actual))
}
