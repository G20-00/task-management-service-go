// Package logger provides a singleton logger instance for application-wide logging.
package logger

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	loggerInstance *logrus.Logger
	once           sync.Once
)

// GetLogger returns the singleton logger instance configured with JSON formatting and file output.
func GetLogger() *logrus.Logger {
	once.Do(func() {
		loggerInstance = logrus.New()

		loggerInstance.SetFormatter(&logrus.JSONFormatter{})

		file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			loggerInstance.SetOutput(os.Stdout)
		} else {
			loggerInstance.SetOutput(io.MultiWriter(os.Stdout, file))
		}

		loggerInstance.SetLevel(logrus.InfoLevel)
	})
	return loggerInstance
}
