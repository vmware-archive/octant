package service

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerHelper(t *testing.T) {
	logger := NewLoggerHelper()
	assert.NotNil(t, logger)
}

func TestSetupPluginLogger(t *testing.T) {
	SetupPluginLogger(Trace)
	assert.Equal(t, log.Prefix(), string(Trace))
	SetupPluginLogger(Debug)
	assert.Equal(t, log.Prefix(), string(Debug))
	SetupPluginLogger(Info)
	assert.Equal(t, log.Prefix(), string(Info))
	SetupPluginLogger(Warn)
	assert.Equal(t, log.Prefix(), string(Warn))
	SetupPluginLogger(Error)
	assert.Equal(t, log.Prefix(), string(Error))
}
