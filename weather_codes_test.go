package openmateo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWeatherCodeDescription(t *testing.T) {
	tests := map[string]struct {
		code int
		want string
	}{
		"Clear sky": {
			code: 0,
			want: "Clear sky",
		},
		"Mainly clear": {
			code: 1,
			want: "Mainly clear",
		},
		"Partly cloudy": {
			code: 2,
			want: "Partly cloudy",
		},
		"Overcast": {
			code: 3,
			want: "Overcast",
		},
		"Fog": {
			code: 45,
			want: "Fog",
		},
		"Depositing rime fog": {
			code: 48,
			want: "Depositing rime fog",
		},
		"Drizzle: Light intensity": {
			code: 51,
			want: "Drizzle: Light intensity",
		},
		"Drizzle: Moderate intensity": {
			code: 53,
			want: "Drizzle: Moderate intensity",
		},
		"Drizzle: Dense intensity": {
			code: 55,
			want: "Drizzle: Dense intensity",
		},
		"Freezing Drizzle: Light intensity": {
			code: 56,
			want: "Freezing Drizzle: Light intensity",
		},
		"Freezing Drizzle: Dense intensity": {
			code: 57,
			want: "Freezing Drizzle: Dense intensity",
		},
		"Rain: Slight intensity": {
			code: 61,
			want: "Rain: Slight intensity",
		},
		"Rain: Moderate intensity": {
			code: 63,
			want: "Rain: Moderate intensity",
		},
		"Rain: Heavy intensity": {
			code: 65,
			want: "Rain: Heavy intensity",
		},
		"Freezing Rain: Light intensity": {
			code: 66,
			want: "Freezing Rain: Light intensity",
		},
		"Freezing Rain: Heavy intensity": {
			code: 67,
			want: "Freezing Rain: Heavy intensity",
		},
		"Snow fall: Slight intensity": {
			code: 71,
			want: "Snow fall: Slight intensity",
		},
		"Snow fall: Moderate intensity": {
			code: 73,
			want: "Snow fall: Moderate intensity",
		},
		"Snow fall: Heavy intensity": {
			code: 75,
			want: "Snow fall: Heavy intensity",
		},
		"Snow grains": {
			code: 77,
			want: "Snow grains",
		},
		"Rain showers: Slight": {
			code: 80,
			want: "Rain showers: Slight",
		},
		"Rain showers: Moderate": {
			code: 81,
			want: "Rain showers: Moderate",
		},
		"Rain showers: Violent": {
			code: 82,
			want: "Rain showers: Violent",
		},
		"Snow showers: Slight": {
			code: 85,
			want: "Snow showers: Slight",
		},
		"Snow showers: Heavy": {
			code: 86,
			want: "Snow showers: Heavy",
		},
		"Thunderstorm: Slight or moderate": {
			code: 95,
			want: "Thunderstorm: Slight or moderate",
		},
		"Thunderstorm with slight hail": {
			code: 96,
			want: "Thunderstorm with slight hail",
		},
		"Thunderstorm with heavy hail": {
			code: 99,
			want: "Thunderstorm with heavy hail",
		},
		"Unknown code": {
			code: 100,
			want: "Unknown code",
		},
		"Negative code": {
			code: -1,
			want: "Unknown code",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := DescribeCode(tc.code)
			assert.Equal(t, tc.want, got)
		})
	}
}
