package provider

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/mrrizkin/omniscan/app/providers/logger/util"
	"github.com/mrrizkin/omniscan/config"
)

type ZeroLogger struct {
	*zerolog.Logger
}

func Zerolog(config *config.App) (*ZeroLogger, error) {
	var writers []io.Writer

	if config.LOG_CONSOLE {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.LOG_FILE {
		rf, err := util.RollingFile(config)
		if err != nil {
			return nil, err
		}
		writers = append(writers, rf)
	}
	mw := io.MultiWriter(writers...)

	switch config.LOG_LEVEL {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "disable":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	zlogInstance := zerolog.New(mw).With().Timestamp().Logger()
	logger := ZeroLogger{&zlogInstance}

	logger.Info(
		"logging configured",
		"fileLogging", config.LOG_FILE,
		"jsonLogOutput", config.LOG_JSON,
		"logDirectory", config.LOG_DIR,
		"fileName", config.NAME+".log",
		"maxSizeMB", config.LOG_MAX_SIZE,
		"maxBackups", config.LOG_MAX_BACKUP,
		"maxAgeInDays", config.LOG_MAX_AGE,
	)

	return &logger, nil
}

// usage
func (z *ZeroLogger) Info(msg string, args ...interface{}) {
	z.argsParser(z.Logger.Info(), args...).Msg(msg)
}

func (z *ZeroLogger) Warn(msg string, args ...interface{}) {
	z.argsParser(z.Logger.Warn(), args...).Msg(msg)
}

func (z *ZeroLogger) Error(msg string, args ...interface{}) {
	z.argsParser(z.Logger.Error(), args...).Msg(msg)
}

func (z *ZeroLogger) Fatal(msg string, args ...interface{}) {
	z.argsParser(z.Logger.Fatal(), args...).Msg(msg)
}

func (z *ZeroLogger) GetLogger() interface{} {
	return z.Logger
}

func (z *ZeroLogger) argsParser(event *zerolog.Event, args ...interface{}) *zerolog.Event {
	if len(args)%2 != 0 {
		z.Warn("logger: args don't match key val")
		return event
	}

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			z.Warn("info: non-string key provided")
			continue
		}

		switch value := args[i+1].(type) {
		case bool:
			event.Bool(key, value)
		case []bool:
			event.Bools(key, value)
		case string:
			event.Str(key, value)
		case []string:
			event.Strs(key, value)
		case int:
			event.Int(key, value)
		case []int:
			event.Ints(key, value)
		case int64:
			event.Int64(key, value)
		case []int64:
			event.Ints64(key, value)
		case float32:
			event.Float32(key, value)
		case float64:
			event.Float64(key, value)
		case time.Time:
			event.Time(key, value)
		case time.Duration:
			event.Dur(key, value)
		case []byte:
			event.Bytes(key, value)
		case error:
			event.Err(value)
		default:
			event.Interface(key, value)
		}
	}

	return event
}
