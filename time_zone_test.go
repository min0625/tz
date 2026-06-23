package tz_test

import (
	"database/sql/driver"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/min0625/tz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustLoadTimeZone(t *testing.T, name string) tz.TimeZone {
	t.Helper()

	z, err := tz.LoadTimeZone(name)
	require.NoError(t, err)

	return z
}

func TestLoadTimeZone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		testName string
		name     string
		want     tz.TimeZone
		wantErr  bool
	}{
		{
			testName: "Empty",
			name:     "",
			want:     tz.TimeZone{},
			wantErr:  false,
		},
		{
			testName: "UTC",
			name:     "UTC",
			want:     tz.TimeZone{},
			wantErr:  false,
		},
		{
			testName: "Local",
			name:     "Local",
			want:     tz.TimeZone{},
			wantErr:  true,
		},
		{
			testName: "America/New_York",
			name:     "America/New_York",
			want:     mustLoadTimeZone(t, "America/New_York"),
			wantErr:  false,
		},
		{
			testName: "InvalidName",
			name:     "Invalid/Zone",
			want:     tz.TimeZone{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			got, err := tz.LoadTimeZone(tt.name)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMustLoadTimeZone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		testName  string
		name      string
		want      tz.TimeZone
		wantPanic bool
	}{
		{
			testName:  "Empty",
			name:      "",
			want:      tz.TimeZone{},
			wantPanic: false,
		},
		{
			testName:  "UTC",
			name:      "UTC",
			want:      tz.TimeZone{},
			wantPanic: false,
		},
		{
			testName:  "America/New_York",
			name:      "America/New_York",
			want:      mustLoadTimeZone(t, "America/New_York"),
			wantPanic: false,
		},
		{
			testName:  "Local",
			name:      "Local",
			wantPanic: true,
		},
		{
			testName:  "InvalidName",
			name:      "Invalid/Zone",
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			if tt.wantPanic {
				assert.Panics(t, func() {
					tz.MustLoadTimeZone(tt.name)
				})

				return
			}

			assert.Equal(t, tt.want, tz.MustLoadTimeZone(tt.name))
		})
	}
}

func TestTimeZone_Location_ZeroValueReturnUTC(t *testing.T) {
	t.Parallel()
	assert.Same(t, tz.TimeZone{}.Location(), time.UTC)
}

func TestTimeZone_Location_NonUTC(t *testing.T) {
	t.Parallel()

	z := mustLoadTimeZone(t, "America/New_York")
	loc := z.Location()

	assert.Equal(t, "America/New_York", loc.String())
	assert.NotSame(t, loc, time.UTC)
}

func TestTimeZone_LoadString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveTimeZone tz.TimeZone
		data         string
		wantTimeZone tz.TimeZone
		wantErr      bool
	}{
		{
			name:         "Empty",
			data:         "",
			wantTimeZone: tz.TimeZone{},
			wantErr:      false,
		},
		{
			name:         "UTC",
			data:         "UTC",
			wantTimeZone: mustLoadTimeZone(t, "UTC"),
			wantErr:      false,
		},
		{
			name:         "America/New_York",
			data:         "America/New_York",
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      false,
		},
		{
			name:         "Asia/Tokyo",
			data:         "Asia/Tokyo",
			wantTimeZone: mustLoadTimeZone(t, "Asia/Tokyo"),
			wantErr:      false,
		},
		{
			name:         "ErrName",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			data:         "ErrName",
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
		{
			name:         "Local",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			data:         "Local",
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.giveTimeZone

			err := z.LoadString(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantTimeZone, z)
		})
	}
}

func TestTimeZone_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		z    tz.TimeZone
		want string
	}{
		{
			name: "Empty",
			z:    mustLoadTimeZone(t, ""),
			want: "UTC",
		},
		{
			name: "UTC",
			z:    mustLoadTimeZone(t, "UTC"),
			want: "UTC",
		},
		{
			name: "America/New_York",
			z:    mustLoadTimeZone(t, "America/New_York"),
			want: "America/New_York",
		},
		{
			name: "Asia/Tokyo",
			z:    mustLoadTimeZone(t, "Asia/Tokyo"),
			want: "Asia/Tokyo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.z.String())
		})
	}
}

func TestTimeZone_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveTimeZone tz.TimeZone
		value        any
		wantTimeZone tz.TimeZone
		wantErr      bool
	}{
		{
			name:         "string_UTC",
			value:        "UTC",
			wantTimeZone: mustLoadTimeZone(t, "UTC"),
			wantErr:      false,
		},
		{
			name:         "string_America/New_York",
			value:        "America/New_York",
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      false,
		},
		{
			name:         "string_Asia/Tokyo",
			value:        "Asia/Tokyo",
			wantTimeZone: mustLoadTimeZone(t, "Asia/Tokyo"),
			wantErr:      false,
		},
		{
			name:         "bytes_America/New_York",
			value:        []byte("America/New_York"),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      false,
		},
		{
			name:         "bytes_ErrName",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			value:        []byte("ErrName"),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
		{
			name:         "string_ErrName",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			value:        "ErrName",
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
		{
			name:         "nil",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			value:        nil,
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
		{
			name:         "unsupported_type",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			value:        []int{1},
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.giveTimeZone

			err := z.Scan(tt.value)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantTimeZone, z)
		})
	}
}

func TestTimeZone_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		z       tz.TimeZone
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "UTC",
			z:       mustLoadTimeZone(t, "UTC"),
			want:    "UTC",
			wantErr: false,
		},
		{
			name:    "Asia/Tokyo",
			z:       mustLoadTimeZone(t, "Asia/Tokyo"),
			want:    "Asia/Tokyo",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.z.Value()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeZone_MarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		z       tz.TimeZone
		want    []byte
		wantErr bool
	}{
		{
			name:    "UTC",
			z:       mustLoadTimeZone(t, "UTC"),
			want:    []byte("UTC"),
			wantErr: false,
		},
		{
			name:    "Asia/Tokyo",
			z:       mustLoadTimeZone(t, "Asia/Tokyo"),
			want:    []byte("Asia/Tokyo"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.z.MarshalText()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeZone_UnmarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveTimeZone tz.TimeZone
		data         []byte
		wantTimeZone tz.TimeZone
		wantErr      bool
	}{
		{
			name:         "UTC",
			data:         []byte("UTC"),
			wantTimeZone: mustLoadTimeZone(t, "UTC"),
			wantErr:      false,
		},
		{
			name:         "America/New_York",
			data:         []byte("America/New_York"),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      false,
		},
		{
			name:         "Asia/Tokyo",
			data:         []byte("Asia/Tokyo"),
			wantTimeZone: mustLoadTimeZone(t, "Asia/Tokyo"),
			wantErr:      false,
		},
		{
			name:         "Empty",
			data:         []byte(""),
			wantTimeZone: tz.TimeZone{},
			wantErr:      false,
		},
		{
			name:         "ErrName",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			data:         []byte("ErrName"),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
		{
			name:         "Local",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			data:         []byte("Local"),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.giveTimeZone

			err := z.UnmarshalText(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantTimeZone, z)
		})
	}
}

func TestTimeZone_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		z       tz.TimeZone
		want    []byte
		wantErr bool
	}{
		{
			name:    "UTC",
			z:       mustLoadTimeZone(t, "UTC"),
			want:    []byte(`"UTC"`),
			wantErr: false,
		},
		{
			name:    "Asia/Tokyo",
			z:       mustLoadTimeZone(t, "Asia/Tokyo"),
			want:    []byte(`"Asia/Tokyo"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.z.MarshalJSON()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeZone_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveTimeZone tz.TimeZone
		data         []byte
		wantTimeZone tz.TimeZone
		wantErr      bool
	}{
		{
			name:         "UTC",
			data:         []byte(`"UTC"`),
			wantTimeZone: mustLoadTimeZone(t, "UTC"),
			wantErr:      false,
		},
		{
			name:         "America/New_York",
			data:         []byte(`"America/New_York"`),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      false,
		},
		{
			name:         "Asia/Tokyo",
			data:         []byte(`"Asia/Tokyo"`),
			wantTimeZone: mustLoadTimeZone(t, "Asia/Tokyo"),
			wantErr:      false,
		},
		{
			name:         "ErrName",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			data:         []byte(`"ErrName"`),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
		{
			name:         "null",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			data:         []byte(`null`),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      false,
		},
		{
			name:         "Local",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			data:         []byte(`"Local"`),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
		{
			name:         "NumberNotString",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			data:         []byte(`123`),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.giveTimeZone

			err := z.UnmarshalJSON(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantTimeZone, z)
		})
	}
}

func TestTimeZone_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"", "UTC", "America/New_York", "Asia/Tokyo", "Europe/London"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			original := mustLoadTimeZone(t, name)

			data, err := original.MarshalJSON()
			require.NoError(t, err)

			var decoded tz.TimeZone
			require.NoError(t, decoded.UnmarshalJSON(data))
			assert.Equal(t, original, decoded)
		})
	}
}

func TestTimeZone_TextRoundTrip(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"", "UTC", "America/New_York", "Asia/Tokyo", "Europe/London"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			original := mustLoadTimeZone(t, name)

			text, err := original.MarshalText()
			require.NoError(t, err)

			var decoded tz.TimeZone
			require.NoError(t, decoded.UnmarshalText(text))
			assert.Equal(t, original, decoded)
		})
	}
}

func TestTimeZone_SQLRoundTrip(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"UTC", "America/New_York", "Asia/Tokyo", "Europe/London"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			original := mustLoadTimeZone(t, name)

			v, err := original.Value()
			require.NoError(t, err)

			var decoded tz.TimeZone
			require.NoError(t, decoded.Scan(v))
			assert.Equal(t, original, decoded)
		})
	}
}
