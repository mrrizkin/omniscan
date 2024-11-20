package util

import (
	"io"
	"os"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/mrrizkin/omniscan/config"
)

func RollingFile(c *config.App) (io.Writer, error) {
	err := os.MkdirAll(c.LOG_DIR, 0744)
	if err != nil {
		return nil, err
	}

	return &lumberjack.Logger{
		Filename:   path.Join(c.LOG_DIR, c.NAME+".log"),
		MaxBackups: c.LOG_MAX_BACKUP, // files
		MaxSize:    c.LOG_MAX_SIZE,   // megabytes
		MaxAge:     c.LOG_MAX_AGE,    // days
	}, nil
}
