package log

import (
	"github.com/sirupsen/logrus"
)

var Logger logrus.FieldLogger = logrus.New().WithField("mod","rmq")

func SetLogger(logger logrus.FieldLogger) {
	Logger = logger
}
