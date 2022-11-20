package tz

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// TimeZone based on time.LoadLocation format, Not support the "Local" time zone.
// The format is "UTC" or IANA time zone database name.
// See: https://www.iana.org/time-zones.
// The zero value means UTC time zone.
// The UTC time zone always uses a zero value.
type TimeZone struct {
	loc *time.Location
}

var UTCTimeZone = TimeZone{}

var (
	_ fmt.Stringer             = TimeZone{}
	_ sql.Scanner              = &TimeZone{}
	_ driver.Valuer            = TimeZone{}
	_ encoding.TextMarshaler   = TimeZone{}
	_ encoding.TextUnmarshaler = &TimeZone{}
	_ json.Marshaler           = TimeZone{}
	_ json.Unmarshaler         = &TimeZone{}
)

func LoadTimeZone(name string) (TimeZone, error) {
	var z TimeZone
	if err := z.loadString(name); err != nil {
		return TimeZone{}, err
	}

	return z, nil
}

// Location always returns a non-nil location.
func (z TimeZone) Location() *time.Location {
	if loc := z.loc; loc != nil {
		return loc
	}

	return time.UTC
}

func (z *TimeZone) loadString(s string) error {
	loc, err := time.LoadLocation(s)
	if err != nil {
		return err
	}

	if loc == time.Local {
		return errors.New("invalid TimeZone: Local")
	}

	if loc == time.UTC {
		loc = nil
	}

	z.loc = loc
	return nil
}

func (z TimeZone) String() string {
	return z.Location().String()
}

func (z *TimeZone) Scan(value any) error {
	var ns sql.NullString
	if err := ns.Scan(value); err != nil {
		return err
	}

	if !ns.Valid {
		return errors.New("converting NULL to TimeZone is unsupported")
	}

	return z.loadString(ns.String)
}

func (z TimeZone) Value() (driver.Value, error) {
	return z.String(), nil
}

func (z TimeZone) MarshalText() (text []byte, err error) {
	return []byte(z.String()), nil
}

func (z *TimeZone) UnmarshalText(text []byte) error {
	return z.loadString(string(text))
}

func (z TimeZone) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", z.String())), nil
}

func (z *TimeZone) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	// See: https://pkg.go.dev/encoding/json#Unmarshaler.
	if string(data) == "null" {
		return nil
	}

	var unquote string
	if _, err := fmt.Sscanf(string(data), `%q`, &unquote); err != nil {
		return err
	}

	return z.loadString(unquote)
}
