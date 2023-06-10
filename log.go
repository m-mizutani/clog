package clog

import (
	"runtime"
	"time"

	"golang.org/x/exp/slog"
)

type Log struct {
	logLevel slog.Level

	// Timestamp is a time when the log is recorded. Format can be specified by WithTimeFmt.
	Timestamp string

	// Elapsed is a time elapsed from the start of the program.
	Elapsed float64

	// Level is a log level. It is one of "DEBUG", "INFO", "WARN", "ERROR", "FATAL".
	Level string

	// Message is a log message.
	Message string

	// FileName is a file name of the source code that calls logger. It is empty if WithSource is not specified.
	FileName string

	// FilePath is a full file path of the source code that calls logger. It is empty if WithSource is not specified.
	FilePath string

	// FuncName is a function name of the source code that calls logger. It is empty if WithSource is not specified.
	FuncName string

	// FuncLine is a line number of the source code that calls logger. It is empty if WithSource is not specified.
	FuncLine int
}

func (x *Log) Coloring(colors *ColorSet) *Log {
	if colors == nil {
		return x
	}

	x.Level = colors.Level[slog.Level(x.logLevel)].SprintFunc()(x.Level)
	x.Timestamp = colors.Time.SprintFunc()(x.Timestamp)
	x.Message = colors.Message.SprintFunc()(x.Message)

	return x
}

var initTime = time.Now()

func elapsedDuration() float64 {
	return time.Since(initTime).Seconds()
}

type source struct {
	FilePath string
	Func     string
	Line     int
}

func getSource(pc uintptr) *source {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()

	return &source{
		FilePath: f.File,
		Func:     f.Function,
		Line:     f.Line,
	}
}
