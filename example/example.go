package main

import (
	"fmt"
	"time"

	"log/slog"

	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	logrusLogger := logrus.New()

	logger := slog.New(sloglogrus.Option{Level: slog.LevelDebug, Logger: logrusLogger}.NewLogrusHandler())
	logger = logger.With("release", "v1.0.0")

	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now().AddDate(0, 0, -1)),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")
}
