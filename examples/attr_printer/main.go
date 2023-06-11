package main

import (
	"os"

	"github.com/m-mizutani/clog"
	"golang.org/x/exp/slog"
)

type User struct {
	Name  string
	Email string
}

type Store struct {
	Name    string
	Address string
	Phone   string
}

func main() {
	user := User{
		Name:  "mizutani",
		Email: "mizutani@hey.com",
	}

	store := Store{
		Name:    "Jiro",
		Address: "Tokyo",
		Phone:   "123-456-7890",
	}
	group := slog.Group("info", slog.Any("user", user), slog.Any("store", store))

	println()
	// This is a default slog handler for comparison
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	slog.New(textHandler).Info("by slog.TextHandler", group)
	println()

	linearHandler := clog.New(clog.WithPrinter(clog.LinearPrinter))
	slog.New(linearHandler).Info("by LinearPrinter", group)
	println()

	prettyHandler := clog.New(clog.WithPrinter(clog.PrettyPrinter))
	slog.New(prettyHandler).Info("by PrettyPrinter", group)
	println()

	indentHandler := clog.New(clog.WithPrinter(clog.IndentPrinter))
	slog.New(indentHandler).Info("by IndentHandler", group)
	println()
}
