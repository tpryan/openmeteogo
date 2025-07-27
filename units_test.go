package openmateo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemperatureUnitString(t *testing.T) {

	tests := map[string]struct {
		input TemperatureUnit
		want  string
	}{
		"celsius": {
			input: Celsius,
			want:  "celsius",
		},
		"fahrenheit": {
			input: Fahrenheit,
			want:  "fahrenheit",
		},
		"unknown": {
			input: 99,
			want:  "unknown",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPrecipitationUnitString(t *testing.T) {

	tests := map[string]struct {
		input PrecipitationUnit
		want  string
	}{
		"in": {
			input: IN,
			want:  "in",
		},
		"mm": {
			input: MM,
			want:  "mm",
		},
		"unknown": {
			input: 99,
			want:  "unknown",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestWindSpeedUnitString(t *testing.T) {
	tests := map[string]struct {
		input WindSpeedUnit
		want  string
	}{
		"kmh": {
			input: KMH,
			want:  "kmh",
		},
		"ms": {
			input: MS,
			want:  "ms",
		},
		"mph": {
			input: MPH,
			want:  "mph",
		},
		"kn": {
			input: KN,
			want:  "kn",
		},
		"unknown": {
			input: 99,
			want:  "unknown",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.String()
			assert.Equal(t, tc.want, got)
		})
	}

}
