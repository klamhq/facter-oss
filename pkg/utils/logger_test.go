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
func TestDefaultLoggerFactory_ImplementsLoggerFactory(t *testing.T) {
	var factory interface{} = &DefaultLoggerFactory{}
	_, ok := factory.(LoggerFactory)
	assert.True(t, ok, "DefaultLoggerFactory should implement LoggerFactory interface")
}

func TestStdLogger_DifferentLevels(t *testing.T) {
	factory := &DefaultLoggerFactory{}
	levels := []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
	for _, level := range levels {
		logger := factory.New(level)
		assert.Equal(t, level, logger.GetLevel(), "Logger level should be set correctly")
	}
}
