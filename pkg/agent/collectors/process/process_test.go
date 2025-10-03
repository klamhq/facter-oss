package process

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProcesses(t *testing.T) {
	logger := logrus.New()
	proc, err := Processes(logger)
	assert.NoError(t, err)
	assert.NotEmpty(t, proc)
}
