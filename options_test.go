package openmateo

import (
	"testing"
	"time"

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

func TestOptionsBuilder(t *testing.T) {
	t.Parallel()

	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	tz, _ := time.LoadLocation("UTC")
	hourly := &HourlyMetrics{Temperature2m: true}
	daily := &DailyMetrics{WeatherCode: true}
	current := &CurrentMetrics{IsDay: true}

	opts := NewOptionsBuilder().
		Latitude(12.34).
		Longitude(56.78).
		TemperatureUnit(Fahrenheit).
		WindspeedUnit(MPH).
		PrecipitationUnit(IN).
		Timezone(*tz).
		PastDays(5).
		ForcastDays(3).
		Start(start).
		End(end).
		HourlyMetrics(hourly).
		DailyMetrics(daily).
		CurrentMetrics(current).
		Build()

	assert.Equal(t, 12.34, opts.Latitude)
	assert.Equal(t, 56.78, opts.Longitude)
	assert.Equal(t, Fahrenheit, opts.TemperatureUnit)
	assert.Equal(t, MPH, opts.WindspeedUnit)
	assert.Equal(t, IN, opts.PrecipitationUnit)
	assert.Equal(t, *tz, opts.Timezone)
	assert.Equal(t, 5, opts.PastDays)
	assert.Equal(t, 3, opts.ForcastDays)
	assert.Equal(t, start, opts.Start)
	assert.Equal(t, end, opts.End)
	assert.Equal(t, hourly, opts.HourlyMetrics)
	assert.Equal(t, daily, opts.DailyMetrics)
	assert.Equal(t, current, opts.CurrentMetrics)
}
