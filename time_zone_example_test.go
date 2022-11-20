package tz_test

import (
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
	fmt.Println(time.Time{}.In(z.Location()).Location().String())

	// Output:
	// America/New_York
	// America/New_York
}

func ExampleTimeZone_zeroValue() {
	z, err := tz.LoadTimeZone("UTC")
	if err != nil {
		panic(err)
	}

	fmt.Println(z == tz.TimeZone{})

	// Output:
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

	v, err := z.Value()
	fmt.Printf("%v %T %v\n", v, v, err)

	// Output:
	// UTC string <nil>
}
