package sloglogrus

import (
	"log/slog"

	"github.com/sirupsen/logrus"
)

var levelMap = map[slog.Level]logrus.Level{
	slog.LevelDebug: logrus.DebugLevel,
	slog.LevelInfo:  logrus.InfoLevel,
	slog.LevelWarn:  logrus.WarnLevel,
	slog.LevelError: logrus.ErrorLevel,
}
