package logger

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	instance *logrus.Logger
	once     sync.Once
)

func GetLogger() *logrus.Logger {
	once.Do(func() {
		instance = logrus.New()

		instance.SetFormatter(&logrus.JSONFormatter{})

		file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			instance.SetOutput(os.Stdout)
		} else {
			instance.SetOutput(io.MultiWriter(os.Stdout, file))
		}

		instance.SetLevel(logrus.InfoLevel)
	})
	return instance
}
