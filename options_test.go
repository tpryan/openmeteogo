// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openmeteogo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {

	tests := map[string]struct {
		input Metrics
		want  string
	}{
		"hourly": {
			input: Metrics{Temperature2m, RelativeHumidity2m},
			want:  "temperature_2m,relative_humidity_2m",
		},
		"daily": {
			input: Metrics{WeatherCode, Temperature2mMax},
			want:  "weather_code,temperature_2m_max",
		},
		"current": {
			input: Metrics{Temperature2m, RelativeHumidity2m},
			want:  "temperature_2m,relative_humidity_2m",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.encode()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestOptionsBuilder(t *testing.T) {

	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	tz, _ := time.LoadLocation("UTC")

	hourly, err := NewMetrics("hourly", Temperature2m, RelativeHumidity2m)
	assert.NoError(t, err)

	daily, err := NewMetrics("daily", WeatherCode)
	assert.NoError(t, err)

	current, err := NewMetrics("current", IsDay)
	assert.NoError(t, err)

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
