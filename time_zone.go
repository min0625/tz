// Package tz provides a TimeZone type based on IANA time zone database names.
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

// TimeZone represents an IANA time zone based on the time.LoadLocation format.
// The "Local" time zone is not supported.
// Valid values are "UTC" or an IANA time zone database name (e.g. "America/New_York").
// See: https://www.iana.org/time-zones.
// The zero value represents the UTC time zone.
// Loading "UTC" or an empty string always results in the zero value.
type TimeZone struct {
	loc *time.Location
}

// UTCTimeZone is the zero value TimeZone representing the UTC time zone.
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

// LoadTimeZone loads a TimeZone by IANA time zone name.
// An empty string and "UTC" both return the zero value.
// Returns an error if the name is invalid or equals "Local".
func LoadTimeZone(name string) (TimeZone, error) {
	var z TimeZone
	if err := z.LoadString(name); err != nil {
		return TimeZone{}, err
	}

	return z, nil
}

// Location returns the *time.Location for this TimeZone.
// The zero value (UTC) returns time.UTC; it never returns nil.
func (z TimeZone) Location() *time.Location {
	if loc := z.loc; loc != nil {
		return loc
	}

	return time.UTC
}

// LoadString loads a TimeZone by IANA time zone name.
// An empty string and "UTC" both set z to the zero value.
// Returns an error if the name is invalid or equals "Local".
// On error z is not modified.
func (z *TimeZone) LoadString(s string) error {
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

// String returns the IANA time zone name, or "UTC" for the zero value.
// It implements the fmt.Stringer interface.
func (z TimeZone) String() string {
	return z.Location().String()
}

// Scan implements the sql.Scanner interface.
// It accepts a string or []byte value containing an IANA time zone name.
// NULL values are not supported and will return an error.
func (z *TimeZone) Scan(value any) error {
	var ns sql.NullString
	if err := ns.Scan(value); err != nil {
		return err
	}

	if !ns.Valid {
		return errors.New("converting NULL to TimeZone is unsupported")
	}

	return z.LoadString(ns.String)
}

// Value implements the driver.Valuer interface.
// It always returns the IANA time zone name as a string driver.Value.
func (z TimeZone) Value() (driver.Value, error) {
	return z.String(), nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// It encodes the time zone as its IANA name (e.g. "America/New_York" or "UTC").
func (z TimeZone) MarshalText() (text []byte, err error) {
	return []byte(z.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It decodes an IANA time zone name from text.
func (z *TimeZone) UnmarshalText(text []byte) error {
	return z.LoadString(string(text))
}

// MarshalJSON implements the json.Marshaler interface.
// It encodes the time zone as a JSON string containing the IANA name.
func (z TimeZone) MarshalJSON() ([]byte, error) {
	return json.Marshal(z.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It decodes an IANA time zone name from a JSON string.
// A JSON null value is ignored and leaves the receiver unchanged.
func (z *TimeZone) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	// See: https://pkg.go.dev/encoding/json#Unmarshaler.
	if string(data) == "null" {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return z.LoadString(s)
}
