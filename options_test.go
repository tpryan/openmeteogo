package openmateo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHourlyMetrics(t *testing.T) {

	tests := map[string]struct {
		input HourlyMetrics
		want  string
	}{
		"basic": {
			input: HourlyMetrics{Temperature2m: true, RelativeHumidity2m: true},
			want:  "temperature_2m,relative_humidity_2m",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.input.encode()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDailyMetrics(t *testing.T) {

	tests := map[string]struct {
		input DailyMetrics
		want  string
	}{
		"basic": {
			input: DailyMetrics{WeatherCode: true, Temperature2mMax: true},
			want:  "weather_code,temperature_2m_max",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.input.encode()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCurrentMetrics(t *testing.T) {

	tests := map[string]struct {
		input CurrentMetrics
		want  string
	}{
		"basic": {
			input: CurrentMetrics{Temperature2m: true, RelativeHumidity2m: true},
			want:  "temperature_2m,relative_humidity_2m",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.input.encode()
			assert.Equal(t, tc.want, got)
		})
	}
}
