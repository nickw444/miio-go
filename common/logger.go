package common

import "github.com/sirupsen/logrus"

var (
	Log logrus.FieldLogger = logrus.New()
)

func SetLogger(logger logrus.FieldLogger) {
	Log = logger
}
