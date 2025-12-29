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
	"fmt"
	"slices"
	"strings"
	"time"
)

// Options holds all the parameters for a weather data request.
// It is configured using the OptionsBuilder.
type Options struct {
	// Latitude for the location to which the weather forecast refers.
	Latitude float64
	// Longitude for the location to which the weather forecast refers.
	Longitude float64
	// TemperatureUnit sets the unit for temperature values. Default is Celsius.
	TemperatureUnit TemperatureUnit
	// WindspeedUnit sets the unit for wind speed values. Default is km/h.
	WindspeedUnit WindSpeedUnit
	// PrecipitationUnit sets the unit for precipitation values. Default is millimeters.
	PrecipitationUnit PrecipitationUnit
	// Timezone for the forecast data. Default is UTC.
	Timezone time.Location
	// PastDays specifies how many days of historical data to retrieve. Default is 0.
	PastDays int
	// ForcastDays specifies how many days of forecast data to retrieve.
	ForcastDays int
	// Start date for the historical data query.
	Start time.Time
	// End date for the historical data query.
	End time.Time
	// Models specifies the weather models to use (e.g. "ecmwf_seas5").
	Models []string
	// HourlyMetrics specifies which hourly weather variables to retrieve.
	HourlyMetrics Metrics
	// DailyMetrics specifies which daily weather variables to retrieve.
	DailyMetrics Metrics
	// WeeklyMetrics specifies which weekly weather variables to retrieve (Seasonal API).
	WeeklyMetrics Metrics
	// MonthlyMetrics specifies which monthly weather variables to retrieve (Seasonal API).
	MonthlyMetrics Metrics
	// CurrentMetrics specifies which current weather variables to retrieve.
	CurrentMetrics Metrics
	// Seasonal forces the request to use the seasonal API endpoint.
	Seasonal bool
}

// OptionsBuilder provides a fluent interface for constructing an Options object.
type OptionsBuilder struct {
	options *Options
}

// NewOptionsBuilder creates a new OptionsBuilder with default options.
func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{options: &Options{}}
}

// Latitude sets the geographical latitude for the request.
func (b *OptionsBuilder) Latitude(lat float64) *OptionsBuilder {
	b.options.Latitude = lat
	return b
}

// Longitude sets the geographical longitude for the request.
func (b *OptionsBuilder) Longitude(lon float64) *OptionsBuilder {
	b.options.Longitude = lon
	return b
}

// TemperatureUnit sets the desired unit for temperature measurements.
func (b *OptionsBuilder) TemperatureUnit(unit TemperatureUnit) *OptionsBuilder {
	b.options.TemperatureUnit = unit
	return b
}

// WindspeedUnit sets the desired unit for wind speed measurements.
func (b *OptionsBuilder) WindspeedUnit(unit WindSpeedUnit) *OptionsBuilder {
	b.options.WindspeedUnit = unit
	return b
}

// PrecipitationUnit sets the desired unit for precipitation measurements.
func (b *OptionsBuilder) PrecipitationUnit(unit PrecipitationUnit) *OptionsBuilder {
	b.options.PrecipitationUnit = unit
	return b
}

// Timezone sets the timezone for the returned data.
func (b *OptionsBuilder) Timezone(tz time.Location) *OptionsBuilder {
	b.options.Timezone = tz
	return b
}

// PastDays sets the number of past days to retrieve data for.
func (b *OptionsBuilder) PastDays(days int) *OptionsBuilder {
	b.options.PastDays = days
	return b
}

// ForcastDays sets the number of future days to retrieve data for.
func (b *OptionsBuilder) ForcastDays(days int) *OptionsBuilder {
	b.options.ForcastDays = days
	return b
}

// Start sets the start date for a specific time-range query.
func (b *OptionsBuilder) Start(start time.Time) *OptionsBuilder {
	b.options.Start = start
	return b
}

// End sets the end date for a specific time-range query.
func (b *OptionsBuilder) End(end time.Time) *OptionsBuilder {
	b.options.End = end
	return b
}

// Models sets the specific weather models to be used.
func (b *OptionsBuilder) Models(models []string) *OptionsBuilder {
	b.options.Models = models
	return b
}

// HourlyMetrics sets the specific hourly metrics to be fetched.
func (b *OptionsBuilder) HourlyMetrics(metrics Metrics) *OptionsBuilder {
	b.options.HourlyMetrics = metrics
	return b
}

// DailyMetrics sets the specific daily metrics to be fetched.
func (b *OptionsBuilder) DailyMetrics(metrics Metrics) *OptionsBuilder {
	b.options.DailyMetrics = metrics
	return b
}

// WeeklyMetrics sets the specific weekly metrics to be fetched (Seasonal API).
func (b *OptionsBuilder) WeeklyMetrics(metrics Metrics) *OptionsBuilder {
	b.options.WeeklyMetrics = metrics
	return b
}

// MonthlyMetrics sets the specific monthly metrics to be fetched (Seasonal API).
func (b *OptionsBuilder) MonthlyMetrics(metrics Metrics) *OptionsBuilder {
	b.options.MonthlyMetrics = metrics
	return b
}

// CurrentMetrics sets the specific current weather metrics to be fetched.
func (b *OptionsBuilder) CurrentMetrics(metrics Metrics) *OptionsBuilder {
	b.options.CurrentMetrics = metrics
	return b
}

// Seasonal forces the request to use the seasonal API endpoint.
func (b *OptionsBuilder) Seasonal(seasonal bool) *OptionsBuilder {
	b.options.Seasonal = seasonal
	return b
}

// Build finalizes the construction and returns the configured Options object.
func (b *OptionsBuilder) Build() *Options {
	return b.options
}

type Metric string

type Metrics []Metric

func (m Metrics) encode() string {
	sl := []string{}
	for _, metric := range m {
		sl = append(sl, string(metric))
	}

	return strings.Join([]string(sl), ",")
}

const (
	Temperature2m               Metric = "temperature_2m"
	RelativeHumidity2m          Metric = "relative_humidity_2m"
	DewPoint2m                  Metric = "dew_point_2m"
	ApparentTemperature         Metric = "apparent_temperature"
	PrecipitationProbability    Metric = "precipitation_probability"
	Precipitation               Metric = "precipitation"
	Rain                        Metric = "rain"
	Showers                     Metric = "showers"
	Snowfall                    Metric = "snowfall"
	SnowDepth                   Metric = "snow_depth"
	WeatherCode                 Metric = "weather_code"
	PressureMsl                 Metric = "pressure_msl"
	SurfacePressure             Metric = "surface_pressure"
	CloudCover                  Metric = "cloud_cover"
	CloudCoverLow               Metric = "cloud_cover_low"
	CloudCoverMid               Metric = "cloud_cover_mid"
	CloudCoverHigh              Metric = "cloud_cover_high"
	Evapotranspiration          Metric = "evapotranspiration"
	Visibility                  Metric = "visibility"
	Et0FaoEvapotranspiration    Metric = "et0_fao_evapotranspiration"
	VapourPressureDeficit       Metric = "vapour_pressure_deficit"
	WindSpeed10m                Metric = "wind_speed_10m"
	WindSpeed80m                Metric = "wind_speed_80m"
	WindSpeed120m               Metric = "wind_speed_120m"
	WindSpeed180m               Metric = "wind_speed_180m"
	WindDirection10m            Metric = "wind_direction_10m"
	WindDirection80m            Metric = "wind_direction_80m"
	WindDirection120m           Metric = "wind_direction_120m"
	WindDirection180m           Metric = "wind_direction_180m"
	WindGusts10m                Metric = "wind_gusts_10m"
	Temperature80m              Metric = "temperature_80m"
	Temperature120m             Metric = "temperature_120m"
	Temperature180m             Metric = "temperature_180m"
	SoilTemperature0cm          Metric = "soil_temperature_0cm"
	SoilTemperature6cm          Metric = "soil_temperature_6cm"
	SoilTemperature18cm         Metric = "soil_temperature_18cm"
	SoilTemperature54cm         Metric = "soil_temperature_54cm"
	SoilMoisture0To1cm          Metric = "soil_moisture_0_to_1cm"
	SoilMoisture1To3cm          Metric = "soil_moisture_1_to_3cm"
	SoilMoisture9To27cm         Metric = "soil_moisture_9_to_27cm"
	SoilMoisture3To9cm          Metric = "soil_moisture_3_to_9cm"
	IsDay                       Metric = "is_day"
	Temperature2mMax            Metric = "temperature_2m_max"
	Temperature2mMin            Metric = "temperature_2m_min"
	ApparentTemperatureMax      Metric = "apparent_temperature_max"
	ApparentTemperatureMin      Metric = "apparent_temperature_min"
	Sunrise                     Metric = "sunrise"
	Sunset                      Metric = "sunset"
	SunshineDuration            Metric = "sunshine_duration"
	DaylightDuration            Metric = "daylight_duration"
	UvIndexMax                  Metric = "uv_index_max"
	UvIndexClearSkyMax          Metric = "uv_index_clear_sky_max"
	RainSum                     Metric = "rain_sum"
	ShowersSum                  Metric = "showers_sum"
	SnowfallSum                 Metric = "snowfall_sum"
	PrecipitationSum            Metric = "precipitation_sum"
	PrecipitationHours          Metric = "precipitation_hours"
	PrecipitationProbabilityMax Metric = "precipitation_probability_max"
	WindSpeed10mMax             Metric = "wind_speed_10m_max"
	WindGusts10mMax             Metric = "wind_gusts_10m_max"
	WindDirection10mDominant    Metric = "wind_direction_10m_dominant"
	ShortwaveRadiationSum       Metric = "shortwave_radiation_sum"

	// Seasonal Metrics (Weekly & Monthly)
	Temperature2mMean          Metric = "temperature_2m_mean"
	Temperature2mAnomaly       Metric = "temperature_2m_anomaly"
	Temperature2mMaxMean       Metric = "temperature_2m_max_mean"
	Temperature2mMinMean       Metric = "temperature_2m_min_mean"
	DewPoint2mMean             Metric = "dew_point_2m_mean"
	PrecipitationMean          Metric = "precipitation_mean"
	PrecipitationAnomaly       Metric = "precipitation_anomaly"
	PressureMslMean            Metric = "pressure_msl_mean"
	PressureMslAnomaly         Metric = "pressure_msl_anomaly"
	SoilMoisture0To10cmMean    Metric = "soil_moisture_0_to_10cm_mean"
	SoilMoisture0To10cmAnomaly Metric = "soil_moisture_0_to_10cm_anomaly"
)

var hourlyMetrics = []Metric{
	Temperature2m,
	RelativeHumidity2m,
	DewPoint2m,
	ApparentTemperature,
	PrecipitationProbability,
	Precipitation,
	Rain,
	Showers,
	Snowfall,
	SnowDepth,
	WeatherCode,
	PressureMsl,
	SurfacePressure,
	CloudCover,
	CloudCoverLow,
	CloudCoverMid,
	CloudCoverHigh,
	Evapotranspiration,
	Visibility,
	Et0FaoEvapotranspiration,
	VapourPressureDeficit,
	WindSpeed10m,
	WindSpeed80m,
	WindSpeed120m,
	WindSpeed180m,
	WindDirection10m,
	WindDirection80m,
	WindDirection120m,
	WindDirection180m,
	WindGusts10m,
	Temperature80m,
	Temperature120m,
	Temperature180m,
	SoilTemperature0cm,
	SoilTemperature6cm,
	SoilTemperature18cm,
	SoilTemperature54cm,
	SoilMoisture0To1cm,
	SoilMoisture1To3cm,
	SoilMoisture9To27cm,
	SoilMoisture3To9cm,
}

func NewMetrics(metricType string, Metrics ...Metric) (Metrics, error) {
	result := []Metric{}

	var allowed []Metric
	switch metricType {
	case "hourly":
		allowed = hourlyMetrics
	case "daily":
		allowed = dailyMetrics
	case "current":
		allowed = currentMetrics
	case "weekly":
		// TODO: Define strict list for weekly if needed
		return Metrics, nil
	case "monthly":
		// TODO: Define strict list for monthly if needed
		return Metrics, nil
	}

	for _, metric := range Metrics {

		if !slices.Contains(allowed, metric) {
			return nil, fmt.Errorf("invalid for %s metrics: %s ", metricType, metric)
		}
		result = append(result, metric)
	}

	return result, nil

}

var currentMetrics = []Metric{
	Temperature2m,
	RelativeHumidity2m,
	IsDay,
	ApparentTemperature,
	Precipitation,
	Rain,
	Showers,
	Snowfall,
	WeatherCode,
	CloudCover,
	PressureMsl,
	SurfacePressure,
	WindSpeed10m,
	WindDirection10m,
	WindGusts10m,
}

var dailyMetrics = []Metric{
	WeatherCode,
	Temperature2mMax,
	Temperature2mMin,
	ApparentTemperatureMax,
	ApparentTemperatureMin,
	Sunrise,
	Sunset,
	SunshineDuration,
	DaylightDuration,
	UvIndexMax,
	UvIndexClearSkyMax,
	RainSum,
	ShowersSum,
	SnowfallSum,
	PrecipitationSum,
	PrecipitationHours,
	PrecipitationProbabilityMax,
	WindSpeed10mMax,
	WindGusts10mMax,
	WindDirection10mDominant,
	ShortwaveRadiationSum,
	Et0FaoEvapotranspiration,
}
