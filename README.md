# clog

Customizable slog.Handler for console.

```go
package main

import (
	"github.com/m-mizutani/clog"
	"log/slog"
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
```

<img width="997" alt="Screenshot 2023-06-11 at 10 41 29" src="https://github.com/m-mizutani/clog/assets/605953/2dcc46ac-2113-44e4-90d1-218bd5fcc12b">

## Options

- `WithWriter`: Output writer. Default is `os.Stdout`.
- `WithLevel`: Log level. Default is `slog.LevelInfo`.
- `WithTimeFmt`: Time format string. Default is `15:04:05.000`.
- `WithColor`: Enable colorized output. Default will be changed by terminal's color support.
- `WithColorMap`: Color map for each log level. Default is `clog.DefaultColorMap`. See [ColorMap](#colormap) section for more detail.
- `WithSource`: Enable source code location. Default is false.
- `WithReplaceAttr`: Replace attribute value. It's same with `slog.ReplaceAttr` in `slog.HandlerOptions`.
- `WithTemplate`: Template string. See [Template](#template) section for more detail.
- `WithAttrPrinter`: Attribute printer. Default is `clog.LinearPrinter`. See [AttrPrinter](#attrprinter) section for more detail.

### ColorMap

You can customize color map for each handler with `clog.ColorMap`. Default is `clog.DefaultColorMap`. If the fields is nil or not set, default color will be used.

- `LogLevel`: You can set color for each log level. If not set, default color `LogLevelDefault` will be used.
- `LogLevelDefault`: Default color for log level.
- `Time`: Color for time string.
- `Message`: Color for log message string.
- `AttrKey`: Color for attribute key string. It's applied or not depends on AttrPrinter.
- `AttrValue`: Color for attribute value string. It's applied or not depends on AttrPrinter.

### Template

Template can be used to customize log format. A developer can use following variables in template string.

- `.Time`: Time string. Format is specified `WithTimeFmt`.
- `.Elapsed`: Duration from the start of the program
- `.Level`: Log level string. e.g. `INFO`, `WARN`, `ERROR`
- `.Message`: Log message
- `.FileName`: A file name of the source code that calls logger. It is empty if WithSource is not specified
- `.FilePath`: A full file path of the source code that calls logger. It is empty if WithSource is not specified.
- `.FileLine`: A line number of the source code that calls logger. It is empty if WithSource is not specified
- `.FuncName` A function name of the source code that calls logger. It is empty if WithSource is not specified

Default is `clog.DefaultTemplate`.

### AttrPrinter

`AttrPrinter` is an interface designed for customizing the way attributes are printed. By default, `clog.LinearPrinter` is used.

- `LinearPrinter`: Print attributes in a linear way.
- `PrettyPrinter`: Print attributes with [pp](https://github.com/k0kubun/pp) package.
- `IndentPrinter`: Print attributes with indent like YAML format.

Full example is [here](./examples/attr_printer/main.go).

```go
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

	linearHandler := clog.New(clog.WithPrinter(clog.LinearPrinter))
	slog.New(linearHandler).Info("by LinearPrinter", group)

	prettyHandler := clog.New(clog.WithPrinter(clog.PrettyPrinter))
	slog.New(prettyHandler).Info("by PrettyPrinter", group)

	indentHandler := clog.New(clog.WithPrinter(clog.IndentPrinter))
	slog.New(indentHandler).Info("by IndentHandler", group)
}
```

<img width="1188" alt="Screenshot 2023-06-11 at 10 39 26" src="https://github.com/m-mizutani/clog/assets/605953/b184644f-080b-41a9-8e5f-16a80d019311">

## License

Apache License 2.0
