package clog

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"log/slog"
)

type ColorMap struct {
	Level        map[slog.Level]*color.Color
	LevelDefault *color.Color

	Time    *color.Color
	Message *color.Color

	// Whether AttrKey and AttrValue color settings are used or not depends on the AttrPrinter
	AttrKey   *color.Color
	AttrValue *color.Color
}

var (
	defaultColorMap    *ColorMap
	enableColorDefault = false
)

func init() {
	defaultColorMap = &ColorMap{
		Level: map[slog.Level]*color.Color{
			slog.LevelDebug: color.New(color.FgWhite, color.Bold),
			slog.LevelInfo:  color.New(color.FgCyan, color.Bold),
			slog.LevelWarn:  color.New(color.FgYellow, color.Bold),
			slog.LevelError: color.New(color.FgRed, color.Bold),
		},
		LevelDefault: color.New(color.FgBlue, color.Bold),
		Time:         color.New(color.FgWhite),
		Message:      color.New(color.FgHiWhite),

		AttrKey:   color.New(color.FgWhite),
		AttrValue: color.New(color.FgHiWhite),
	}

	colorTerminals := []string{
		"xterm",
		"vt100",
		"rxvt",
		"screen",
	}
	if v, ok := os.LookupEnv("TERM"); ok {
		for _, t := range colorTerminals {
			if strings.Contains(v, t) {
				enableColorDefault = true
				break
			}
		}
	}
}
