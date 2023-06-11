package main

import (
	"github.com/m-mizutani/clog"
	"golang.org/x/exp/slog"
)

func main() {
	handler := clog.New(
		clog.WithColor(true),
		clog.WithSource(true),
	)
	logger := slog.New(handler)

	logger.Info("hello, world!", slog.String("foo", "bar"))
	logger.Warn("What?", slog.Group("group1", slog.String("foo", "bar")))
	logger.WithGroup("hex").Error("Ouch!", slog.Int("num", 123))
}
