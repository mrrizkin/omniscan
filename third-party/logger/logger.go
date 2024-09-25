package logger

import (
	"github.com/mrrizkin/omniscan/system/config"
)

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(err error, msg string)
	Fatal(err error, msg string)
	GetLogger() interface{}
}

func Zerolog(config *config.Config) (Logger, error) {
	return newZerolog(config)
}
