package clog

import (
	"bytes"
	"io"
	"os"
	"text/template"

	"golang.org/x/exp/slog"
)

// config is the configuration for the handler. The struct is immutable after creation.
type config struct {
	w           io.Writer
	level       slog.Leveler
	timeFmt     string
	addSource   bool
	enableColor bool
	replaceAttr func(groups []string, a slog.Attr) slog.Attr
	printer     func(w io.Writer) AttrPrinter
	colors      *ColorSet
	tmpl        *template.Template
}

func newConfig() *config {
	return &config{
		w:           os.Stdout,
		level:       slog.LevelInfo,
		timeFmt:     "2006-01-02 15:04:05",
		addSource:   false,
		enableColor: enableColorDefault,
		printer:     LinearPrinter,

		colors: defaultColorSet,
		tmpl:   defaultTmpl,
	}
}

const DefaultTemplate = `{{.Elapsed | printf "%8.3f" }} {{.Level}} {{ if .FileName }}[{{.FileName}}:{{.FuncLine}}] {{ end }}{{.Message}} `

var defaultTmpl *template.Template

func init() {
	tmpl, err := template.New("default").Parse(DefaultTemplate)
	if err != nil {
		panic(err)
	}
	defaultTmpl = tmpl
}

type Option func(*config)

// WithWriter sets the writer for the handler. The default is os.Stdout.
func WithWriter(w io.Writer) Option {
	return func(cfg *config) {
		cfg.w = w
	}
}

// WithLevel sets the minimum level for the handler. The default is LevelInfo.
func WithLevel(level slog.Leveler) Option {
	return func(cfg *config) {
		cfg.level = level
	}
}

// WithTimeFmt sets the time format for the time attribute. The default is "2006-01-02 15:04:05".
func WithTimeFmt(timeFmt string) Option {
	return func(cfg *config) {
		cfg.timeFmt = timeFmt
	}
}

// WithColor enables or disables color output. The default is enabled.
func WithColor(color bool) Option {
	return func(cfg *config) {
		cfg.enableColor = color
	}
}

// WithSource enables or disables adding the source attribute. The default is disable.
func WithSource(addSource bool) Option {
	return func(cfg *config) {
		cfg.addSource = addSource
	}
}

// WithReplaceAttr sets the function for replacing attributes. The default is nil.
func WithReplaceAttr(replaceAttr func(groups []string, a slog.Attr) slog.Attr) Option {
	return func(cfg *config) {
		cfg.replaceAttr = replaceAttr
	}
}

// WithPrinter sets the printer for printing attributes. The default is LinearPrinter.
func WithPrinter(printer func(w io.Writer) AttrPrinter) Option {
	return func(cfg *config) {
		cfg.printer = printer
	}
}

// WithColorSet sets the color set for the handler.
func WithColorSet(colors *ColorSet) Option {
	return func(cfg *config) {
		cfg.colors = colors
	}
}

// WithTemplate sets the template for the handler. The default is DefaultTemplate. This option executes dry run and panics if the template is invalid.
func WithTemplate(tmpl *template.Template) Option {
	return func(cfg *config) {
		// dry run
		log := &Log{
			Timestamp: "2006-01-02 15:04:05",
			Elapsed:   1.23456789,
			Level:     "INFO",
			Message:   "hello, world!",
			FileName:  "foo.go",
			FilePath:  "/path/to/foo.go",
			FuncName:  "main",
			FuncLine:  10,
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, log); err != nil {
			panic(err)
		}

		cfg.tmpl = tmpl
	}
}
