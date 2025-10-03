package utils

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCheckBinInstalled_BinExists(t *testing.T) {
	logger := logrus.New()
	result := CheckBinInstalled(logger, "go")
	assert.True(t, result)
}

func TestCheckBinInstalled_BinDoesNotExist(t *testing.T) {
	logger := logrus.New()
	result := CheckBinInstalled(logger, "nonexistentbinary123")
	assert.False(t, result)
}
