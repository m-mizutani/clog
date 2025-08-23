# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`clog` is a customizable slog.Handler for console output in Go. It provides colored, formatted console logging with configurable templates and attribute printers.

## Restrictions and Rules

### Directory

- When you are mentioned about `tmp` directory, you SHOULD NOT see `/tmp`. You need to check `./tmp` directory from root of the repository.

### Exposure policy

In principle, do not trust developers who use this library from outside

- Do not export unnecessary methods, structs, and variables
- Assume that exposed items will be changed. Never expose fields that would be problematic if changed
- Use `export_test.go` for items that need to be exposed for testing purposes
- **Exception**: Domain models (`pkg/domain/model/*`) can have exported fields as they represent data structures

### Firestore Struct Tags

- **NEVER use firestore struct tags on domain models**
- Domain models should be pure Go structs without any persistence-specific annotations
- This keeps the domain layer independent of the infrastructure layer

### Check

When making changes, before finishing the task, always:
- Run `go vet ./...`, `go fmt ./...` to format the code
- Run `golangci-lint run ./...` to check lint error
- Run `gosec -exclude-generated -quiet ./...` to check security issue
- Run `go test ./...` to check side effect
- **For GraphQL changes: Run `task graphql` and verify no compilation errors**
- **For GraphQL changes: Check frontend GraphQL queries are updated accordingly**

### Language

All comment and character literal in source code must be in English

### Testing

- Test files should have `package {name}_test`. Do not use same package name
- **🚨 CRITICAL RULE: Test MUST be included in same name test file. (e.g. test for `abc.go` must be in `abc_test.go`) 🚨**

#### Repository Testing Strategy
- **🚨 CRITICAL: Repository tests MUST be placed in `pkg/repository/database/` directory with common test suites**
- Create shared test functions that verify identical behavior across all repository implementations (Firestore, Memory, etc.)
- Each repository implementation must pass the exact same test suite to ensure behavioral consistency
- Use a common test interface pattern to test all implementations uniformly
- This ensures that switching between repository implementations (e.g., Memory for testing, Firestore for production) maintains identical behavior

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