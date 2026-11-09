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

// Marine support.
//
// The marine wave model is a global *ocean* model. A request for a position
// outside its water mask - a marina well up an estuary, for example - does not
// fail: the API answers HTTP 200 with a full-length series of JSON nulls.
//
// Decoding those nulls into []float64 would turn them into zeros, and a caller
// checking len(resp.Daily.WaveHeightMax) > 0 would read a dead-flat calm sea
// where there is in fact no data at all. That is a silent correctness problem
// in a safety-adjacent path, so every marine value is a *float64: nil means
// "no value here", 0 means a real zero reading.
//
// Use HasAny to tell a covered position from an out-of-domain one.

// marinePath is the path of the marine endpoint. The host it is served from
// is Client.marineHost.
const marinePath = "/v1/marine"

// Value returns the first non-nil entry and whether one existed.
func Value(vals []*float64) (float64, bool) {
	for _, v := range vals {
		if v != nil {
			return *v, true
		}
	}

	return 0, false
}

// HasAny reports whether any entry is non-nil. False across a marine series
// means the position lies outside the wave model's domain.
func HasAny(vals []*float64) bool {
	_, ok := Value(vals)

	return ok
}

// Floats flattens a pointer slice, substituting zero for nil entries. Only safe
// once HasAny is true, where nils are gaps inside a covered series rather than
// a wholesale absence of data.
func Floats(vals []*float64) []float64 {
	if vals == nil {
		return nil
	}

	out := make([]float64, len(vals))
	for i, v := range vals {
		if v != nil {
			out[i] = *v
		}
	}

	return out
}

// HourlyMarine holds slices for each hourly marine metric. It is embedded in
// Hourly, so the fields are reachable as hourly.WaveHeight and friends.
//
// A nil entry means the API returned null for that hour, which for a marine
// metric usually means the position lies outside the wave model's domain. A
// non-nil entry pointing at 0 is a real reading of zero.
type HourlyMarine struct {
	WaveHeight                  []*float64 `json:"wave_height"`
	WaveDirection               []*float64 `json:"wave_direction"`
	WavePeriod                  []*float64 `json:"wave_period"`
	WavePeakPeriod              []*float64 `json:"wave_peak_period"`
	WindWaveHeight              []*float64 `json:"wind_wave_height"`
	WindWaveDirection           []*float64 `json:"wind_wave_direction"`
	WindWavePeriod              []*float64 `json:"wind_wave_period"`
	WindWavePeakPeriod          []*float64 `json:"wind_wave_peak_period"`
	SwellWaveHeight             []*float64 `json:"swell_wave_height"`
	SwellWaveDirection          []*float64 `json:"swell_wave_direction"`
	SwellWavePeriod             []*float64 `json:"swell_wave_period"`
	SwellWavePeakPeriod         []*float64 `json:"swell_wave_peak_period"`
	SecondarySwellWaveHeight    []*float64 `json:"secondary_swell_wave_height"`
	SecondarySwellWaveDirection []*float64 `json:"secondary_swell_wave_direction"`
	SecondarySwellWavePeriod    []*float64 `json:"secondary_swell_wave_period"`
	TertiarySwellWaveHeight     []*float64 `json:"tertiary_swell_wave_height"`
	TertiarySwellWaveDirection  []*float64 `json:"tertiary_swell_wave_direction"`
	TertiarySwellWavePeriod     []*float64 `json:"tertiary_swell_wave_period"`
	SeaLevelHeight              []*float64 `json:"sea_level_height_msl"`
	SeaSurfaceTemperature       []*float64 `json:"sea_surface_temperature"`
	OceanCurrentVelocity        []*float64 `json:"ocean_current_velocity"`
	OceanCurrentDirection       []*float64 `json:"ocean_current_direction"`
}

// HourlyMarineUnits describes the units for the hourly marine data. It is
// embedded in HourlyUnits.
type HourlyMarineUnits struct {
	WaveHeight                  string `json:"wave_height"`
	WaveDirection               string `json:"wave_direction"`
	WavePeriod                  string `json:"wave_period"`
	WavePeakPeriod              string `json:"wave_peak_period"`
	WindWaveHeight              string `json:"wind_wave_height"`
	WindWaveDirection           string `json:"wind_wave_direction"`
	WindWavePeriod              string `json:"wind_wave_period"`
	WindWavePeakPeriod          string `json:"wind_wave_peak_period"`
	SwellWaveHeight             string `json:"swell_wave_height"`
	SwellWaveDirection          string `json:"swell_wave_direction"`
	SwellWavePeriod             string `json:"swell_wave_period"`
	SwellWavePeakPeriod         string `json:"swell_wave_peak_period"`
	SecondarySwellWaveHeight    string `json:"secondary_swell_wave_height"`
	SecondarySwellWaveDirection string `json:"secondary_swell_wave_direction"`
	SecondarySwellWavePeriod    string `json:"secondary_swell_wave_period"`
	TertiarySwellWaveHeight     string `json:"tertiary_swell_wave_height"`
	TertiarySwellWaveDirection  string `json:"tertiary_swell_wave_direction"`
	TertiarySwellWavePeriod     string `json:"tertiary_swell_wave_period"`
	SeaLevelHeight              string `json:"sea_level_height_msl"`
	SeaSurfaceTemperature       string `json:"sea_surface_temperature"`
	OceanCurrentVelocity        string `json:"ocean_current_velocity"`
	OceanCurrentDirection       string `json:"ocean_current_direction"`
}

// DailyMarine holds slices for each daily marine metric. It is embedded in
// Daily, so the fields are reachable as daily.WaveHeightMax and friends.
//
// A nil entry means the API returned null for that day. See HourlyMarine for
// what that implies.
type DailyMarine struct {
	WaveHeightMax              []*float64 `json:"wave_height_max"`
	WaveDirectionDominant      []*float64 `json:"wave_direction_dominant"`
	WavePeriodMax              []*float64 `json:"wave_period_max"`
	WindWaveHeightMax          []*float64 `json:"wind_wave_height_max"`
	WindWaveDirectionDominant  []*float64 `json:"wind_wave_direction_dominant"`
	WindWavePeriodMax          []*float64 `json:"wind_wave_period_max"`
	WindWavePeakPeriodMax      []*float64 `json:"wind_wave_peak_period_max"`
	SwellWaveHeightMax         []*float64 `json:"swell_wave_height_max"`
	SwellWaveDirectionDominant []*float64 `json:"swell_wave_direction_dominant"`
	SwellWavePeriodMax         []*float64 `json:"swell_wave_period_max"`
	SwellWavePeakPeriodMax     []*float64 `json:"swell_wave_peak_period_max"`
}

// DailyMarineUnits describes the units for the daily marine data. It is
// embedded in DailyUnits.
type DailyMarineUnits struct {
	WaveHeightMax              string `json:"wave_height_max"`
	WaveDirectionDominant      string `json:"wave_direction_dominant"`
	WavePeriodMax              string `json:"wave_period_max"`
	WindWaveHeightMax          string `json:"wind_wave_height_max"`
	WindWaveDirectionDominant  string `json:"wind_wave_direction_dominant"`
	WindWavePeriodMax          string `json:"wind_wave_period_max"`
	WindWavePeakPeriodMax      string `json:"wind_wave_peak_period_max"`
	SwellWaveHeightMax         string `json:"swell_wave_height_max"`
	SwellWaveDirectionDominant string `json:"swell_wave_direction_dominant"`
	SwellWavePeriodMax         string `json:"swell_wave_period_max"`
	SwellWavePeakPeriodMax     string `json:"swell_wave_peak_period_max"`
}

// CurrentMarine holds the current marine values. It is embedded in Current.
//
// A nil field means the API returned null, not that the sea is flat.
type CurrentMarine struct {
	WaveHeight            *float64 `json:"wave_height"`
	WaveDirection         *float64 `json:"wave_direction"`
	WavePeriod            *float64 `json:"wave_period"`
	WindWaveHeight        *float64 `json:"wind_wave_height"`
	WindWaveDirection     *float64 `json:"wind_wave_direction"`
	WindWavePeriod        *float64 `json:"wind_wave_period"`
	WindWavePeakPeriod    *float64 `json:"wind_wave_peak_period"`
	SwellWaveHeight       *float64 `json:"swell_wave_height"`
	SwellWaveDirection    *float64 `json:"swell_wave_direction"`
	SwellWavePeriod       *float64 `json:"swell_wave_period"`
	SwellWavePeakPeriod   *float64 `json:"swell_wave_peak_period"`
	SeaLevelHeight        *float64 `json:"sea_level_height_msl"`
	SeaSurfaceTemperature *float64 `json:"sea_surface_temperature"`
	OceanCurrentVelocity  *float64 `json:"ocean_current_velocity"`
	OceanCurrentDirection *float64 `json:"ocean_current_direction"`
}

// CurrentMarineUnits describes the units for the current marine data. It is
// embedded in CurrentUnits.
type CurrentMarineUnits struct {
	WaveHeight            string `json:"wave_height"`
	WaveDirection         string `json:"wave_direction"`
	WavePeriod            string `json:"wave_period"`
	WindWaveHeight        string `json:"wind_wave_height"`
	WindWaveDirection     string `json:"wind_wave_direction"`
	WindWavePeriod        string `json:"wind_wave_period"`
	WindWavePeakPeriod    string `json:"wind_wave_peak_period"`
	SwellWaveHeight       string `json:"swell_wave_height"`
	SwellWaveDirection    string `json:"swell_wave_direction"`
	SwellWavePeriod       string `json:"swell_wave_period"`
	SwellWavePeakPeriod   string `json:"swell_wave_peak_period"`
	SeaLevelHeight        string `json:"sea_level_height_msl"`
	SeaSurfaceTemperature string `json:"sea_surface_temperature"`
	OceanCurrentVelocity  string `json:"ocean_current_velocity"`
	OceanCurrentDirection string `json:"ocean_current_direction"`
}

// The marine Metric constants live in options.go, alongside every other
// metric the library exposes.

// marineHourlyMetrics are the metrics the marine endpoint accepts for hourly.
var marineHourlyMetrics = []Metric{
	WaveHeight,
	WaveDirection,
	WavePeriod,
	WavePeakPeriod,
	WindWaveHeight,
	WindWaveDirection,
	WindWavePeriod,
	WindWavePeakPeriod,
	SwellWaveHeight,
	SwellWaveDirection,
	SwellWavePeriod,
	SwellWavePeakPeriod,
	SecondarySwellWaveHeight,
	SecondarySwellWaveDirection,
	SecondarySwellWavePeriod,
	TertiarySwellWaveHeight,
	TertiarySwellWaveDirection,
	TertiarySwellWavePeriod,
	SeaLevelHeight,
	SeaSurfaceTemperature,
	OceanCurrentVelocity,
	OceanCurrentDirection,
}

// marineDailyMetrics are the metrics the marine endpoint accepts for daily.
var marineDailyMetrics = []Metric{
	WaveHeightMax,
	WaveDirectionDominant,
	WavePeriodMax,
	WindWaveHeightMax,
	WindWaveDirectionDominant,
	WindWavePeriodMax,
	WindWavePeakPeriodMax,
	SwellWaveHeightMax,
	SwellWaveDirectionDominant,
	SwellWavePeriodMax,
	SwellWavePeakPeriodMax,
}

// marineCurrentMetrics are the metrics the marine endpoint accepts for current.
var marineCurrentMetrics = []Metric{
	WaveHeight,
	WaveDirection,
	WavePeriod,
	WindWaveHeight,
	WindWaveDirection,
	WindWavePeriod,
	WindWavePeakPeriod,
	SwellWaveHeight,
	SwellWaveDirection,
	SwellWavePeriod,
	SwellWavePeakPeriod,
	SeaLevelHeight,
	SeaSurfaceTemperature,
	OceanCurrentVelocity,
	OceanCurrentDirection,
}
