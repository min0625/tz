# Golang Time Zone Type
[![Go Reference](https://pkg.go.dev/badge/github.com/min0625/tz.svg)](https://pkg.go.dev/github.com/min0625/tz)

## Features
- TimeZone based on time.LoadLocation format, Not support the "Local" time zone.
- The format is "UTC" or IANA time zone database name. See: https://www.iana.org/time-zones.
- The zero value means UTC time zone.
- The UTC time zone always uses a zero value.
- Implement the fmt.Stringer interface.
- Implement the sql.Scanner interface.
- Implement the driver.Valuer interface.
- Implement the encoding.TextMarshaler interface.
- Implement the encoding.TextUnmarshaler interface.
- Implement the json.Marshaler interface.
- Implement the json.Unmarshaler interface.

## Installation
```sh
go get github.com/min0625/tz
```

## Quick start
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
