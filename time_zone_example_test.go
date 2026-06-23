package tz_test

import (
	"encoding/json"
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/min0625/tz"
)

func ExampleTimeZone() {
	z, err := tz.LoadTimeZone("America/New_York")
	if err != nil {
		panic(err)
	}

	fmt.Println(z.String())

	// Location returns a *time.Location that can be applied to a time.Time.
	loc := z.Location()
	fmt.Println(time.Time{}.In(loc).Location().String())

	// Output:
	// America/New_York
	// America/New_York
}

func ExampleTimeZone_zeroValue() {
	// Loading "UTC" (or an empty string) always yields the zero value.
	z, err := tz.LoadTimeZone("UTC")
	if err != nil {
		panic(err)
	}

	fmt.Println(z == tz.TimeZone{})

	// Output:
	// true
}

func ExampleTimeZone_Equal() {
	a := tz.MustLoadTimeZone("America/New_York")
	b := tz.MustLoadTimeZone("America/New_York")

	// Each load returns a distinct *time.Location, so == is unreliable,
	// but Equal compares by IANA name.
	fmt.Println(a == b)
	fmt.Println(a.Equal(b))

	// The zero value equals "UTC".
	fmt.Println(tz.TimeZone{}.Equal(tz.MustLoadTimeZone("UTC")))

	// Output:
	// false
	// true
	// true
}

func ExampleLoadTimeZone_errors() {
	// "Local" is always rejected.
	_, err := tz.LoadTimeZone("Local")
	fmt.Println(err != nil)

	// An unknown name is rejected too.
	_, err = tz.LoadTimeZone("Not/A_Zone")
	fmt.Println(err != nil)

	// Output:
	// true
	// true
}

func ExampleMustLoadTimeZone() {
	z := tz.MustLoadTimeZone("America/New_York")
	fmt.Println(z.String())

	// Output:
	// America/New_York
}

func ExampleTimeZone_LoadString() {
	var z tz.TimeZone
	if err := z.LoadString("Asia/Tokyo"); err != nil {
		panic(err)
	}

	fmt.Println(z)

	// Loading "UTC" (or an empty string) resets z to the zero value.
	if err := z.LoadString("UTC"); err != nil {
		panic(err)
	}

	fmt.Println(z == tz.TimeZone{})

	// Output:
	// Asia/Tokyo
	// true
}

func ExampleTimeZone_Scan() {
	var z tz.TimeZone
	if err := z.Scan("America/New_York"); err != nil {
		panic(err)
	}

	fmt.Println(time.Time{}.In(z.Location()).Location().String())

	// Output:
	// America/New_York
}

func ExampleTimeZone_Value() {
	z, err := tz.LoadTimeZone("America/New_York")
	if err != nil {
		panic(err)
	}

	v, err := z.Value()
	fmt.Printf("%v %T %v\n", v, v, err)

	// Output:
	// America/New_York string <nil>
}

func ExampleTimeZone_Value_zeroValue() {
	var z tz.TimeZone

	// The zero value is stored as the "UTC" string, never as NULL.
	v, err := z.Value()
	fmt.Printf("%v %T %v\n", v, v, err)

	// Output:
	// UTC string <nil>
}

func ExampleUTCTimeZone() {
	z := tz.MustLoadTimeZone("UTC")
	fmt.Println(z == tz.UTCTimeZone)

	// Output:
	// true
}

func ExampleTimeZone_AppendText() {
	z := tz.MustLoadTimeZone("Asia/Tokyo")

	b := []byte("zone=")
	b, _ = z.AppendText(b)
	fmt.Println(string(b))

	// Output:
	// zone=Asia/Tokyo
}

func ExampleTimeZone_MarshalText() {
	z := tz.MustLoadTimeZone("Asia/Tokyo")
	text, _ := z.MarshalText()
	fmt.Println(string(text))

	// Output:
	// Asia/Tokyo
}

func ExampleTimeZone_UnmarshalText() {
	var z tz.TimeZone

	_ = z.UnmarshalText([]byte("Asia/Tokyo"))
	fmt.Println(z)

	// Output:
	// Asia/Tokyo
}

func ExampleTimeZone_MarshalJSON() {
	type Event struct {
		Name     string      `json:"name"`
		TimeZone tz.TimeZone `json:"time_zone"`
	}

	e := Event{
		Name:     "Meeting",
		TimeZone: tz.MustLoadTimeZone("Europe/London"),
	}

	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

	// Output:
	// {"name":"Meeting","time_zone":"Europe/London"}
}

func ExampleTimeZone_UnmarshalJSON() {
	type Event struct {
		Name     string      `json:"name"`
		TimeZone tz.TimeZone `json:"time_zone"`
	}

	data := []byte(`{"name":"Meeting","time_zone":"Europe/London"}`)

	var e Event

	_ = json.Unmarshal(data, &e)
	fmt.Println(e.Name, e.TimeZone)

	// Output:
	// Meeting Europe/London
}
