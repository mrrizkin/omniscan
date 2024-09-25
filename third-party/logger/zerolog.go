package logger

import (
	"io"
	"os"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/mrrizkin/omniscan/system/config"
	log "github.com/rs/zerolog"
)

type zerolog struct {
	*log.Logger
}

func newZerolog(config *config.Config) (Logger, error) {
	var writers []io.Writer

	if config.LOG_CONSOLE {
		writers = append(writers, log.ConsoleWriter{Out: os.Stderr})
	}
	if config.LOG_FILE {
		rf, err := rollingFile(config)
		if err != nil {
			return nil, err
		}
		writers = append(writers, rf)
	}
	mw := io.MultiWriter(writers...)

	switch config.LOG_LEVEL {
	case "panic":
		log.SetGlobalLevel(log.PanicLevel)
	case "fatal":
		log.SetGlobalLevel(log.FatalLevel)
	case "error":
		log.SetGlobalLevel(log.ErrorLevel)
	case "warn":
		log.SetGlobalLevel(log.WarnLevel)
	case "info":
		log.SetGlobalLevel(log.InfoLevel)
	case "debug":
		log.SetGlobalLevel(log.DebugLevel)
	case "trace":
		log.SetGlobalLevel(log.TraceLevel)
	case "disable":
		log.SetGlobalLevel(log.Disabled)
	}

	logger := log.New(mw).With().Timestamp().Logger()

	logger.Info().
		Bool("fileLogging", config.LOG_FILE).
		Bool("jsonLogOutput", config.LOG_JSON).
		Str("logDirectory", config.LOG_DIR).
		Str("fileName", config.APP_NAME+".log").
		Int("maxSizeMB", config.LOG_MAX_SIZE).
		Int("maxBackups", config.LOG_MAX_BACKUP).
		Int("maxAgeInDays", config.LOG_MAX_AGE).
		Msg("logging configured")

	return &zerolog{
		Logger: &logger,
	}, nil
}

func (z *zerolog) Info(msg string) {
	z.Logger.Info().Msg(msg)
}

func (z *zerolog) Warn(msg string) {
	z.Logger.Warn().Msg(msg)
}

func (z *zerolog) Error(err error, msg string) {
	z.Logger.Error().Err(err).Msg(msg)
}

func (z *zerolog) Fatal(err error, msg string) {
	z.Logger.Fatal().Err(err).Msg(msg)
}

func (z *zerolog) GetLogger() interface{} {
	return z.Logger
}

func rollingFile(c *config.Config) (io.Writer, error) {
	err := os.MkdirAll(c.LOG_DIR, 0744)
	if err != nil {
		return nil, err
	}

	return &lumberjack.Logger{
		Filename:   path.Join(c.LOG_DIR, c.APP_NAME+".log"),
		MaxBackups: c.LOG_MAX_BACKUP, // files
		MaxSize:    c.LOG_MAX_SIZE,   // megabytes
		MaxAge:     c.LOG_MAX_AGE,    // days
	}, nil
}
