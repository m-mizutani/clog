package main

import (
	"log/slog"
	"os"

	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/clog/hooks"
	"github.com/m-mizutani/goerr/v2"
)

func someAction(args string) error {
	return goerr.New("something wrong", goerr.V("args", args))
}

func main() {
	options := []clog.Option{
		clog.WithColor(false),
	}

	if _, ok := os.LookupEnv("HANDLE_ERROR"); ok {
		options = append(options, clog.WithAttrHook(hooks.GoErr()))
	}

	logger := slog.New(clog.New(options...))

	err := someAction("foo")
	logger.Error("oops", "error", err)
}
