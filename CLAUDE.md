# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`clog` is a customizable slog.Handler for console output in Go. It provides colored, formatted console logging with configurable templates and attribute printers.

## Commands

### Testing
- Run all tests: `go test ./...`
- Run tests with verbose output: `go test ./... -v`
- Run specific test: `go test -run TestName`

### Code Quality
- Format code: `go fmt ./...`
- Run vet checks: `go vet ./...`

### Build
- Build the module: `go build ./...`
- Run examples: `go run examples/basic/main.go`

## Architecture

### Core Components

The package is structured around the slog.Handler interface with these main components:

1. **Handler** (`handler.go`): Main slog.Handler implementation that manages log output, grouping, and attribute handling.

2. **Config** (`config.go`): Configuration system using functional options pattern. Key options include:
   - Output writer, log level, time format
   - Color enablement and color mapping
   - Source code location tracking
   - Template customization
   - Attribute hooks and printers

3. **AttrPrinter Interface** (`printer.go`): Pluggable system for attribute formatting with three implementations:
   - `LinearPrinter`: Default inline attribute printing
   - `PrettyPrinter`: Uses k0kubun/pp for pretty printing
   - `IndentPrinter`: YAML-like indented format

4. **Color System** (`color.go`): Manages terminal colors using fatih/color, with customizable ColorMap for different log elements.

5. **Attribute Hooks** (`attr.go`): Extension point for custom attribute handling, includes built-in `GoerrHook` for github.com/m-mizutani/goerr/v2 error formatting.

### Key Design Patterns

- **Functional Options**: All configuration through `With*` functions
- **Immutable Config**: Configuration is immutable after handler creation
- **Clone Pattern**: Handler cloning for WithGroup/WithAttrs operations
- **Template System**: Go text/template for customizable log format

### Dependencies

- `github.com/fatih/color`: Terminal color support
- `github.com/k0kubun/pp/v3`: Pretty printing for PrettyPrinter
- `github.com/m-mizutani/goerr/v2`: Enhanced error handling with GoerrHook
- `github.com/m-mizutani/gt`: Testing utilities (dev dependency)