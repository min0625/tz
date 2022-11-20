package tz

import (
	"database/sql/driver"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustLoadTimeLocation(t *testing.T, name string) *time.Location {
	loc, err := time.LoadLocation(name)
	require.NoError(t, err)
	return loc
}

func mustLoadTimeZone(t *testing.T, name string) TimeZone {
	z, err := LoadTimeZone(name)
	require.NoError(t, err)
	return z
}

func TestLoadTimeZone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		testName string
		name     string
		want     TimeZone
		wantErr  bool
	}{
		{
			testName: "Empty",
			name:     "",
			want:     TimeZone{},
			wantErr:  false,
		},
		{
			testName: "UTC",
			name:     "UTC",
			want:     TimeZone{},
			wantErr:  false,
		},
		{
			testName: "Local",
			name:     "Local",
			want:     TimeZone{},
			wantErr:  true,
		},
		{
			testName: "America/New_York",
			name:     "America/New_York",
			want: TimeZone{
				loc: mustLoadTimeLocation(t, "America/New_York"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			got, err := LoadTimeZone(tt.name)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestTimeZone_Location_ZeroValueReturnUTC(t *testing.T) {
	t.Parallel()
	assert.Same(t, TimeZone{}.Location(), time.UTC)
}

func TestTimeZone_loadString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveTimeZone TimeZone
		data         string
		wantTimeZone TimeZone
		wantErr      bool
	}{
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
			name:         "America/New_York",
			data:         "America/New_York",
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      false,
		},
		{
			name:         "ErrName",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			data:         "ErrName",
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.giveTimeZone

			err := z.loadString(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, z, tt.wantTimeZone)
		})
	}
}

func TestTimeZone_String(t *testing.T) {
	tests := []struct {
		name string
		z    TimeZone
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.z.String(), tt.want)
		})
	}
}

func TestTimeZone_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveTimeZone TimeZone
		value        any
		wantTimeZone TimeZone
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
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.giveTimeZone

			err := z.Scan(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, z, tt.wantTimeZone)
		})
	}
}

func TestTimeZone_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		z       TimeZone
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.z.Value()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestTimeZone_MarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		z       TimeZone
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.z.MarshalText()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestTimeZone_UnmarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveTimeZone TimeZone
		data         []byte
		wantTimeZone TimeZone
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
			name:         "America/New_York",
			data:         []byte("America/New_York"),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      false,
		},
		{
			name:         "ErrName",
			giveTimeZone: mustLoadTimeZone(t, "America/New_York"),
			data:         []byte("ErrName"),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.giveTimeZone

			err := z.UnmarshalText(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, z, tt.wantTimeZone)
		})
	}
}

func TestTimeZone_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		z       TimeZone
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.z.MarshalJSON()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestTimeZone_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveTimeZone TimeZone
		data         []byte
		wantTimeZone TimeZone
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
			name:         "America/New_York",
			data:         []byte(`"America/New_York"`),
			wantTimeZone: mustLoadTimeZone(t, "America/New_York"),
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
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.giveTimeZone

			err := z.UnmarshalJSON(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, z, tt.wantTimeZone)
		})
	}
}
