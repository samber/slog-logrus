package sloglogrus

import (
	"log/slog"
	"testing"

	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
)

func TestLogrusHandler_WithGroupSiblings(t *testing.T) {
	logger, hook := logrustest.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	root := slog.New(Option{Logger: logger}.NewLogrusHandler())
	parent := root.WithGroup("a").WithGroup("b").WithGroup("c")

	first := parent.WithGroup("d")
	second := parent.WithGroup("e")

	first.Info("first", "k", 1)
	second.Info("second", "k", 2)

	entries := hook.AllEntries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	for i, group := range []string{"d", "e"} {
		got := map[string]any(entries[i].Data)
		inner, ok := got["a"].(map[string]any)["b"].(map[string]any)["c"].(map[string]any)[group]
		if !ok || inner == nil {
			t.Fatalf("record %d lost group %q: %v", i, group, got)
		}
	}
}
