# Golang Time Zone Type
[![Go Reference](https://pkg.go.dev/badge/github.com/min0625/tz.svg)](https://pkg.go.dev/github.com/min0625/tz)

## Features
- Based on `time.LoadLocation` format; the `"Local"` time zone is not supported.
- Accepts `"UTC"`, an empty string, or any IANA time zone database name (e.g. `"America/New_York"`). See: https://www.iana.org/time-zones.
- The zero value represents the UTC time zone; loading `"UTC"` or `""` always produces the zero value.
- Implements `fmt.Stringer`
- Implements `sql.Scanner`
- Implements `driver.Valuer`
- Implements `encoding.TextMarshaler`
- Implements `encoding.TextUnmarshaler`
- Implements `json.Marshaler`
- Implements `json.Unmarshaler`

## Installation
```sh
go get github.com/min0625/tz
```

## Quick start

> **Note**: Import `_ "time/tzdata"` to embed the IANA time zone database directly into your binary,
> so it works correctly in environments where the system timezone data may be absent (e.g. scratch/Alpine containers).

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

	fmt.Println(z.String())
	fmt.Println(time.Time{}.In(z.Location()).Location().String())

	// Output:
	// America/New_York
	// America/New_York
}
```

## Example
See: [./time_zone_example_test.go](./time_zone_example_test.go)
