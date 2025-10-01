package utils

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// LoggerFactory is an interface for creating loggers with a specific log level
type LoggerFactory interface {
	New(level logrus.Level) *logrus.Logger
}

// DefaultLoggerFactory is the default implementation of LoggerFactory
type DefaultLoggerFactory struct{}

// StdLogger give a configured logger
func (f *DefaultLoggerFactory) New(level logrus.Level) *logrus.Logger {
	logger := &logrus.Logger{}
	logger.SetFormatter(new(prefixed.TextFormatter))
	logger.SetOutput(os.Stdout)
	logger.SetLevel(level)
	logger.SetReportCaller(true)
	return logger
}
