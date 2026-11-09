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
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// f returns a pointer to v, for building expected []*float64 values.
func f(v float64) *float64 {
	return &v
}

// marineJSON builds a marine response with the given hourly wave_height and
// daily wave_height_max fragments.
func marineJSON(hourly, daily string) string {
	return `{
		"hourly": {"time": ["2026-01-01T00:00", "2026-01-01T01:00", "2026-01-01T02:00"], "wave_height": ` + hourly + `},
		"daily": {"time": ["2026-01-01", "2026-01-02"], "wave_height_max": ` + daily + `}
	}`
}

// TestMarineDecode_AllNull covers the Benicia Marina case: a position outside
// the wave model's domain answers HTTP 200 with a full-length series of nulls.
// Decoded into []float64 those would read as a dead-flat calm sea.
//
// TestParsing checks the whole struct against the same fixture; this checks
// the property a caller actually relies on.
func TestMarineDecode_AllNull(t *testing.T) {
	dat, err := os.ReadFile("test_data/marine_out_of_domain.json")
	require.NoError(t, err)

	var wd WeatherData
	require.NoError(t, json.Unmarshal(dat, &wd))

	// The slices keep their length, so a len() check still sees a populated
	// response. Only the nils tell the truth.
	require.Len(t, wd.Hourly.WaveHeight, 3)
	require.Len(t, wd.Daily.WaveHeightMax, 2)

	for i, v := range wd.Hourly.WaveHeight {
		assert.Nil(t, v, "hourly entry %d", i)
	}

	for i, v := range wd.Daily.WaveHeightMax {
		assert.Nil(t, v, "daily entry %d", i)
	}

	assert.False(t, HasAny(wd.Hourly.WaveHeight))
	assert.False(t, HasAny(wd.Daily.WaveHeightMax))
}

// TestMarineDecode_Mixed covers gaps inside an otherwise covered series.
func TestMarineDecode_Mixed(t *testing.T) {
	var wd WeatherData
	err := json.Unmarshal([]byte(marineJSON(`[1.2, null, 0.9]`, `[1.2, null]`)), &wd)
	require.NoError(t, err)

	require.Len(t, wd.Hourly.WaveHeight, 3)
	assert.Equal(t, f(1.2), wd.Hourly.WaveHeight[0])
	assert.Nil(t, wd.Hourly.WaveHeight[1])
	assert.Equal(t, f(0.9), wd.Hourly.WaveHeight[2])

	// A gap does not make the position uncovered.
	assert.True(t, HasAny(wd.Hourly.WaveHeight))
	assert.True(t, HasAny(wd.Daily.WaveHeightMax))
}

// TestMarineDecode_RealZero is the assertion the whole change exists for: a
// measured calm must stay distinguishable from absent data.
func TestMarineDecode_RealZero(t *testing.T) {
	var wd WeatherData
	err := json.Unmarshal([]byte(marineJSON(`[0.0, 0.0, 0.0]`, `[0.0, 0.0]`)), &wd)
	require.NoError(t, err)

	require.Len(t, wd.Hourly.WaveHeight, 3)
	for i, v := range wd.Hourly.WaveHeight {
		require.NotNil(t, v, "hourly entry %d should be a real zero, not absent", i)
		assert.Equal(t, 0.0, *v)
	}

	for i, v := range wd.Daily.WaveHeightMax {
		require.NotNil(t, v, "daily entry %d should be a real zero, not absent", i)
		assert.Equal(t, 0.0, *v)
	}

	assert.True(t, HasAny(wd.Hourly.WaveHeight))
	assert.True(t, HasAny(wd.Daily.WaveHeightMax))
}

// TestMarineDecode_RealZeroVsNull pins the distinction directly: the same
// field, same length, same len() check, opposite meaning.
func TestMarineDecode_RealZeroVsNull(t *testing.T) {
	var calm, absent WeatherData
	require.NoError(t, json.Unmarshal([]byte(marineJSON(`[0.0, 0.0, 0.0]`, `[0.0, 0.0]`)), &calm))
	require.NoError(t, json.Unmarshal([]byte(marineJSON(`[null, null, null]`, `[null, null]`)), &absent))

	// Indistinguishable by length, which is what made the original bug silent.
	require.Len(t, calm.Hourly.WaveHeight, 3)
	require.Len(t, absent.Hourly.WaveHeight, 3)

	// And flattening loses the distinction too, which is why Floats is only
	// safe once HasAny is true.
	assert.Equal(t, Floats(calm.Hourly.WaveHeight), Floats(absent.Hourly.WaveHeight))

	// HasAny is the thing that tells them apart.
	assert.True(t, HasAny(calm.Hourly.WaveHeight))
	assert.False(t, HasAny(absent.Hourly.WaveHeight))
}

// TestMarineDecode_AbsentFields guards the case where the caller asked for no
// marine metrics at all: the fields stay nil rather than becoming empty.
func TestMarineDecode_AbsentFields(t *testing.T) {
	var wd WeatherData
	err := json.Unmarshal([]byte(`{"latitude": 52.52, "hourly": {"time": ["2026-01-01T00:00"], "temperature_2m": [3.4]}}`), &wd)
	require.NoError(t, err)

	assert.Nil(t, wd.Hourly.WaveHeight)
	assert.False(t, HasAny(wd.Hourly.WaveHeight))
	assert.Nil(t, Floats(wd.Hourly.WaveHeight))
	assert.Equal(t, []float64{3.4}, wd.Hourly.Temperature2m)
}

func TestValue(t *testing.T) {
	tests := map[string]struct {
		in     []*float64
		want   float64
		wantOk bool
	}{
		"nil slice":       {in: nil, want: 0, wantOk: false},
		"empty":           {in: []*float64{}, want: 0, wantOk: false},
		"all nil":         {in: []*float64{nil, nil}, want: 0, wantOk: false},
		"first set":       {in: []*float64{f(1.2), nil}, want: 1.2, wantOk: true},
		"leading gap":     {in: []*float64{nil, f(0.9), f(1.1)}, want: 0.9, wantOk: true},
		"real zero first": {in: []*float64{f(0), f(1.1)}, want: 0, wantOk: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := Value(tc.in)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantOk, ok)
		})
	}
}

func TestHasAny(t *testing.T) {
	tests := map[string]struct {
		in   []*float64
		want bool
	}{
		"nil slice":  {in: nil, want: false},
		"empty":      {in: []*float64{}, want: false},
		"all nil":    {in: []*float64{nil, nil, nil}, want: false},
		"one set":    {in: []*float64{nil, f(1.2), nil}, want: true},
		"real zeros": {in: []*float64{f(0), f(0)}, want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, HasAny(tc.in))
		})
	}
}

func TestFloats(t *testing.T) {
	tests := map[string]struct {
		in   []*float64
		want []float64
	}{
		"nil slice": {in: nil, want: nil},
		"empty":     {in: []*float64{}, want: []float64{}},
		"no gaps":   {in: []*float64{f(1.2), f(0.9)}, want: []float64{1.2, 0.9}},
		"with gap":  {in: []*float64{f(1.2), nil, f(0.9)}, want: []float64{1.2, 0, 0.9}},
		"all nil":   {in: []*float64{nil, nil}, want: []float64{0, 0}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, Floats(tc.in))
		})
	}
}

// TestFloats_DoesNotAliasInput guards against Floats handing back a slice that
// shares storage with, or mutates, its input.
func TestFloats_DoesNotAliasInput(t *testing.T) {
	in := []*float64{f(1.2), nil, f(0.9)}
	out := Floats(in)

	out[0] = 99

	require.NotNil(t, in[0])
	assert.Equal(t, 1.2, *in[0])
	assert.Nil(t, in[1])
}

func TestNewMetrics_Marine(t *testing.T) {
	tests := map[string]struct {
		metricType string
		metrics    []Metric
		wantErr    bool
	}{
		"marine hourly": {
			metricType: "marine_hourly",
			metrics:    []Metric{WaveHeight, SwellWavePeriod, OceanCurrentVelocity},
		},
		"marine daily": {
			metricType: "marine_daily",
			metrics:    []Metric{WaveHeightMax, SwellWavePeakPeriodMax},
		},
		"marine current": {
			metricType: "marine_current",
			metrics:    []Metric{WaveHeight, SeaSurfaceTemperature},
		},
		"secondary swell is hourly only": {
			metricType: "marine_current",
			metrics:    []Metric{SecondarySwellWaveHeight},
			wantErr:    true,
		},
		"daily metric rejected on marine hourly": {
			metricType: "marine_hourly",
			metrics:    []Metric{WaveHeightMax},
			wantErr:    true,
		},
		"land metric rejected on marine hourly": {
			metricType: "marine_hourly",
			metrics:    []Metric{Temperature2m},
			wantErr:    true,
		},
		"marine metric rejected on land hourly": {
			metricType: "hourly",
			metrics:    []Metric{WaveHeight},
			wantErr:    true,
		},
		"marine metric rejected on land daily": {
			metricType: "daily",
			metrics:    []Metric{WaveHeightMax},
			wantErr:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewMetrics(tc.metricType, tc.metrics...)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, Metrics(tc.metrics), got)
		})
	}
}
