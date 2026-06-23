# Go Time Zone Type

[![Go Reference](https://pkg.go.dev/badge/github.com/min0625/tz.svg)](https://pkg.go.dev/github.com/min0625/tz)
[![codecov](https://codecov.io/gh/min0625/tz/branch/main/graph/badge.svg)](https://codecov.io/gh/min0625/tz)

**English** | [繁體中文](./README.zh-TW.md)

A small Go library providing a `TimeZone` type backed by the IANA time zone database.

## Features

- **Zero value is UTC** — `TimeZone{}` represents UTC; loading `"UTC"` or `""` always yields the zero value.
- **IANA names** — accepts any [IANA time zone database](https://www.iana.org/time-zones) name (e.g. `"America/New_York"`).
- **No `"Local"`** — the `"Local"` time zone is always rejected with an error.
- **Equality** — compare values with the `Equal` method; `==` is only reliable against the zero value, since a `TimeZone` holds a `*time.Location` pointer.
- **Rich interface support** — implements `fmt.Stringer`, `sql.Scanner`, `driver.Valuer`, `encoding.TextMarshaler/Unmarshaler`, and `json.Marshaler/Unmarshaler`.
- **No dependencies** — built entirely on the standard library.

## Installation

```sh
go get github.com/min0625/tz
```

Requires Go 1.24 or later.

## Quick start

> **Note**: Import `_ "time/tzdata"` to embed the IANA time zone database directly into your binary,
> so it works correctly in environments where system timezone data may be absent (e.g. scratch or Alpine containers).

```go
package main

import (
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/min0625/tz"
)

func main() {
	z, err := tz.LoadTimeZone("America/New_York")
	if err != nil {
		panic(err)
	}

	fmt.Println(z)                                      // America/New_York
	fmt.Println(time.Now().In(z.Location()).Location()) // America/New_York

	// Compare with Equal; == is unreliable for loaded zones.
	other := tz.MustLoadTimeZone("America/New_York")
	fmt.Println(z.Equal(other)) // true

	// UTC is the zero value.
	utc, _ := tz.LoadTimeZone("UTC")
	fmt.Println(utc == tz.TimeZone{}) // true
}
```

## Documentation

- API reference: [pkg.go.dev/github.com/min0625/tz](https://pkg.go.dev/github.com/min0625/tz)
- Runnable examples: [time_zone_example_test.go](./time_zone_example_test.go)

## License

See [LICENSE](./LICENSE).
