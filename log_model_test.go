package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateLog(t *testing.T) {
	log := generateLog()
	assert.NotNil(t, log)
	assert.NotNil(t, log.Content)
	assert.Equal(t, 100, len(log.Content))
}

func TestCalculateLogBytes(t *testing.T) {
	log := generateLog()
	assert.Equal(t, 100, calculateLogBytes(log))
}

func TestCalculateLogsBytes(t *testing.T) {
	logs := []*Log{}
	for i := 0; i < 100; i++ {
		logs = append(logs, generateLog())
	}
	assert.Equal(t, 10000, calculateLogsBytes(logs))
}
