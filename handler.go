package sloglogrus

import (
	"context"

	"log/slog"

	"github.com/sirupsen/logrus"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// optional: logrus logger (default: logrus.StandardLogger())
	Logger *logrus.Logger

	// optional: customize json payload builder
	Converter Converter
}

func (o Option) NewLogrusHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Logger == nil {
		// should be selected lazily ?
		o.Logger = logrus.StandardLogger()
	}

	return &LogrusHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*LogrusHandler)(nil)

type LogrusHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *LogrusHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *LogrusHandler) Handle(ctx context.Context, record slog.Record) error {
	converter := DefaultConverter
	if h.option.Converter != nil {
		converter = h.option.Converter
	}

	level := levelMap[record.Level]
	args := converter(h.attrs, &record)

	logrus.NewEntry(h.option.Logger).
		WithContext(ctx).
		WithTime(record.Time).
		WithFields(logrus.Fields(args)).
		Log(level, record.Message)

	return nil
}

func (h *LogrusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogrusHandler{
		option: h.option,
		attrs:  appendAttrsToGroup(h.groups, h.attrs, attrs),
		groups: h.groups,
	}
}

func (h *LogrusHandler) WithGroup(name string) slog.Handler {
	return &LogrusHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
