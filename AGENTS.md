# Project Guidelines

## Overview

`github.com/min0625/tz` is a small, single-file Go library that provides a `TimeZone` type backed by the IANA time zone database.
It exposes no external dependencies beyond the standard library (test-only dependency: `testify`).

## Architecture

- **`time_zone.go`** — the entire public API: `TimeZone` struct, `LoadTimeZone`, `LoadString`, and implementations of `fmt.Stringer`, `sql.Scanner`, `driver.Valuer`, `encoding.TextMarshaler/Unmarshaler`, `json.Marshaler/Unmarshaler`.
- **`time_zone_test.go`** — table-driven unit tests in package `tz_test`.
- **`time_zone_example_test.go`** — runnable `Example*` tests in package `tz_test`.

## Build and Test

```sh
# Run all tests (with race detector)
make test
# equivalent: go test -v -race -failfast ./...

# Lint (golangci-lint)
make lint

# Lint with auto-fix
make fix

# Lint + test
make check
```

Embed the IANA timezone database in tests by importing `_ "time/tzdata"`.

## Code Conventions

- **Go version**: 1.24+. Use modern idioms (e.g. `any` instead of `interface{}`, no `tt := tt` loop-variable copy).
- **Formatting**: `gofmt`, `gofumpt`, `goimports`, `gci`, `golines` — run `make fix` to auto-apply.
- **Tests**: all tests must be parallel (`t.Parallel()`). Use table-driven tests. Test package is `tz_test` (black-box).
- **UTC invariant**: loading `"UTC"` or `""` must always produce the zero value `TimeZone{}`.
- **No Local**: `"Local"` must always be rejected with an error; `time.Local` must never be stored.
- **Interfaces**: any new serialisation format must maintain the UTC-as-zero-value invariant. SQL NULL must be rejected explicitly (return an error). JSON `null` follows the Go standard convention — ignore it and leave the receiver unchanged (see [`encoding/json.Unmarshaler`](https://pkg.go.dev/encoding/json#Unmarshaler)).

## Linting

golangci-lint is configured in `.golangci.yaml`. Key enabled linters: `errcheck`, `staticcheck`, `gosec`, `errorlint`, `testifylint`, `testpackage`, `wsl_v5`. Do not add `//nolint` directives without a clear justification comment.
