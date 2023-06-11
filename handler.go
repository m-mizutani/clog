package clog

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

// Handler is a slog handler that writes logs to an io.Writer.
type Handler struct {
	cfg *config

	attrs []slog.Attr
	group string
	mutex *sync.Mutex

	parent *Handler
}

var _ slog.Handler = (*Handler)(nil)

// New creates a new handler.
func New(options ...Option) *Handler {
	h := &Handler{
		cfg:   newConfig(),
		mutex: &sync.Mutex{},
	}

	for _, option := range options {
		option(h.cfg)
	}

	return h
}

// clone returns a copy of the handler.
func (x *Handler) clone() *Handler {
	newHandler := &Handler{
		cfg:    x.cfg,
		parent: x,
		mutex:  x.mutex,
	}

	return newHandler
}

// Enabled implements slog.Handler.
func (x *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return x.cfg.level.Level() <= level
}

type stack struct {
	handlers []*Handler
}

func (x *stack) push(h *Handler) {
	x.handlers = append(x.handlers, h)
}

func (x *stack) pop() *Handler {
	if len(x.handlers) == 0 {
		return nil
	}

	h := x.handlers[len(x.handlers)-1]
	x.handlers = x.handlers[:len(x.handlers)-1]
	return h
}

// Handle implements slog.Handler.
func (x *Handler) Handle(ctx context.Context, record slog.Record) error {
	x = x.clone()
	buf := &bytes.Buffer{}

	log := &Log{
		logLevel:  record.Level,
		Timestamp: record.Time.Format(x.cfg.timeFmt),
		Elapsed:   elapsedDuration(),
		Level:     record.Level.String(),
		Message:   record.Message,
	}
	if record.Time.IsZero() {
		log.Timestamp = "(no time)"
	}

	if x.cfg.addSource && record.PC != 0 {
		src := getSource(record.PC)
		log.FileName = filepath.Base(src.FilePath)
		log.FilePath = src.FilePath
		log.FuncName = src.Func
		log.FileLine = src.Line
	}

	if x.cfg.enableColor {
		log = log.Coloring(x.cfg.colors)
	}

	if err := x.cfg.tmpl.Execute(buf, log); err != nil {
		return goerr.Wrap(err, "failed to execute template")
	}

	// print attrs
	record.Attrs(func(attr slog.Attr) bool {
		x.attrs = append(x.attrs, attr)
		return true
	})

	printer := x.cfg.newPrinter(buf, x.cfg)
	st := &stack{}
	for handler := x; handler != nil; handler = handler.parent {
		st.push(handler)
	}

	printHandlerAttrs(printer, st)

	fmt.Fprint(buf, "\n")

	x.mutex.Lock()
	defer x.mutex.Unlock()
	if _, err := x.cfg.w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func printHandlerAttrs(p AttrPrinter, st *stack) {
	h := st.pop()
	if h == nil {
		return
	}

	if h.group != "" {
		p.Enter(h.group)
	}

	printAttrs(p, h.attrs)
	printHandlerAttrs(p, st)

	if h.group != "" {
		p.Exit(h.group)
	}
}

func printAttrs(p AttrPrinter, attrs []slog.Attr) {
	for _, attr := range attrs {
		if attr.Equal(slog.Attr{}) {
			continue // ignored
		}

		switch attr.Value.Kind() {
		case slog.KindGroup:
			p.Enter(attr.Key)
			printAttrs(p, attr.Value.Group())
			p.Exit(attr.Key)

		default:
			p.Print(attr)
		}
	}
}

// WithAttrs implements slog.Handler.
func (x *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := x.clone()
	newHandler.attrs = attrs
	return newHandler
}

// WithGroup implements slog.Handler.
func (x *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return x
	}

	newHandler := x.clone()
	newHandler.group = name
	return newHandler
}
