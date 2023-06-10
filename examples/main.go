package main

import (
	"github.com/m-mizutani/clog"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(clog.New(clog.WithColor(true), clog.WithSource(true)))

	logger.Info("hello, world!", slog.String("foo", "bar"))
	logger.Warn("What?", slog.Group("group1", slog.String("foo", "bar")))
	logger.WithGroup("hex").Error("Ouch!", slog.Int("num", 123))
}