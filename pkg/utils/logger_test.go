package utils

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func TestStdLogger_Level(t *testing.T) {
	factory := &DefaultLoggerFactory{}
	logger := factory.New(logrus.InfoLevel)
	assert.Equal(t, logrus.InfoLevel, logger.GetLevel())
}

func TestStdLogger_Formatter(t *testing.T) {
	factory := &DefaultLoggerFactory{}
	logger := factory.New(logrus.WarnLevel)
	_, ok := logger.Formatter.(*prefixed.TextFormatter)
	assert.True(t, ok, "Logger formatter should be prefixed.TextFormatter")
}

func TestStdLogger_Output(t *testing.T) {
	factory := &DefaultLoggerFactory{}
	logger := factory.New(logrus.ErrorLevel)
	assert.NotNil(t, logger.Out)
}

func TestStdLogger_ReportCaller(t *testing.T) {
	factory := &DefaultLoggerFactory{}
	logger := factory.New(logrus.DebugLevel)
	assert.True(t, logger.ReportCaller)
}
