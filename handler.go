package sloglogrus

import (
	"context"
	"log/slog"

	slogcommon "github.com/samber/slog-common"
	"github.com/sirupsen/logrus"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// optional: logrus logger (default: logrus.StandardLogger())
	Logger *logrus.Logger

	// optional: if set, GetLogger will be used to get the logger at each call to Handle
	GetLogger func(ctx context.Context) *logrus.Logger

	// optional: customize json payload builder
	Converter Converter

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
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

	level := LogLevels[record.Level]
	args := converter(h.option.AddSource, h.option.ReplaceAttr, h.attrs, h.groups, &record)

	logger := h.option.Logger
	if h.option.GetLogger != nil {
		logger = h.option.GetLogger(ctx)
	}

	logrus.NewEntry(logger).
		WithContext(ctx).
		WithTime(record.Time).
		WithFields(logrus.Fields(args)).
		Log(level, record.Message)

	return nil
}

func (h *LogrusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogrusHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
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
