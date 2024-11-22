package logger

import (
	"go.uber.org/fx/fxevent"

	"github.com/mrrizkin/omniscan/app/providers/logger/provider"
	"github.com/mrrizkin/omniscan/config"
)

type LoggerProvider interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	FxLogEvent(event fxevent.Event)
	GetLogger() interface{}
}

type Logger struct {
	provider LoggerProvider
}

func (*Logger) Construct() interface{} {
	return func(config *config.App) (*Logger, error) {
		p, err := provider.Zerolog(config)
		if err != nil {
			return nil, err
		}

		return &Logger{provider: p}, nil
	}
}

// usage
func (log *Logger) Info(msg string, args ...interface{}) {
	log.provider.Info(msg, args...)
}

func (log *Logger) Warn(msg string, args ...interface{}) {
	log.provider.Warn(msg, args...)
}

func (log *Logger) Error(msg string, args ...interface{}) {
	args = append(args, "stack_trace", stackTrace())
	log.provider.Error(msg, args...)
}

func (log *Logger) Fatal(msg string, args ...interface{}) {
	args = append(args, "stack_trace", stackTrace())
	log.provider.Fatal(msg, args...)
}

func (log *Logger) GetLogger() interface{} {
	return log.provider.GetLogger()
}

func (log *Logger) LogEvent(event fxevent.Event) {
	log.provider.FxLogEvent(event)
}
