package logger

import (
	neasted "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"os"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&neasted.Formatter{
		HideKeys: true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func NewLoggerWithFields(field map[string]interface{}) *logrus.Entry {
	return log.WithFields(logrus.Fields(field))
}

func AssertMessage(field map[string]interface{}, message string) {
	log.WithFields(logrus.Fields(field)).Error(message)
	os.Exit(1)
}
