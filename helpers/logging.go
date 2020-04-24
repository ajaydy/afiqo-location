package helpers

import (
	"github.com/sirupsen/logrus"
	"os"
)

type (
	Logger struct {
		Out *logrus.Logger
		Err *logrus.Logger
	}
)

func NewLogger() *Logger {
	return &Logger{
		Out: &logrus.Logger{
			Formatter: new(logrus.TextFormatter),
			Out:       os.Stdout,
			Level:     logrus.InfoLevel,
		},
		Err: &logrus.Logger{
			Formatter: new(logrus.TextFormatter),
			Out:       os.Stdout,
			Level:     logrus.InfoLevel,
		},
	}
}
