The Marine Weather API follows a similar pattern but uses a different host and specific wave/ocean variables.

Here is the plan to add Marine Forecast support:

1. **Update `options.go**`: Add the specific Marine metrics (e.g., `WaveHeight`, `SwellWaveHeight`, `OceanCurrentVelocity`) to the `Metrics` constants.
2. **Update `openmeteo.go**`:
* Add the `marineHost` constant (`marine-api.open-meteo.com`).
* Add a `GetMarine` method to the `Client`.
* Implement a `marineUrl` helper method.
* Extend the `Hourly`, `HourlyUnits`, `Daily`, and `DailyUnits` structs to include the new marine data fields.



### **Implementation**

#### **1. `options.go**`

I have added the Marine-specific metrics for both Hourly and Daily intervals.

```go
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
	// Models specifies the weather models to use (e.g. "ecmwf_seas5", "ecmwf_wam025").
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
	return strings.Join(sl, ",")
}

// NewMetrics creates a new Metrics slice, validating that the provided metrics are valid for the given type.
// Valid types are "hourly", "daily", "weekly", "monthly", and "current".
func NewMetrics(kind string, metrics ...Metric) (Metrics, error) {
	// For now, we will just return the metrics as is.
	// In the future, we can add validation logic here.
	if !slices.Contains([]string{"hourly", "daily", "weekly", "monthly", "current"}, kind) {
		return nil, fmt.Errorf("invalid metric kind: %s", kind)
	}
	return metrics, nil
}

const (
	// Hourly Metrics
	Temperature2m            Metric = "temperature_2m"
	RelativeHumidity2m       Metric = "relative_humidity_2m"
	DewPoint2m               Metric = "dew_point_2m"
	ApparentTemperature      Metric = "apparent_temperature"
	PrecipitationProbability Metric = "precipitation_probability"
	Precipitation            Metric = "precipitation"
	Rain                     Metric = "rain"
	Showers                  Metric = "showers"
	Snowfall                 Metric = "snowfall"
	SnowDepth                Metric = "snow_depth"
	WeatherCode              Metric = "weather_code"
	PressureMsl              Metric = "pressure_msl"
	SurfacePressure          Metric = "surface_pressure"
	CloudCover               Metric = "cloud_cover"
	CloudCoverLow            Metric = "cloud_cover_low"
	CloudCoverMid            Metric = "cloud_cover_mid"
	CloudCoverHigh           Metric = "cloud_cover_high"
	Evapotranspiration       Metric = "evapotranspiration"
	Visibility               Metric = "visibility"
	Et0FaoEvapotranspiration Metric = "et0_fao_evapotranspiration"
	VapourPressureDeficit    Metric = "vapour_pressure_deficit"
	WindSpeed10m             Metric = "wind_speed_10m"
	WindSpeed80m             Metric = "wind_speed_80m"
	WindSpeed120m            Metric = "wind_speed_120m"
	WindSpeed180m            Metric = "wind_speed_180m"
	WindDirection10m         Metric = "wind_direction_10m"
	WindDirection80m         Metric = "wind_direction_80m"
	WindDirection120m        Metric = "wind_direction_120m"
	WindDirection180m        Metric = "wind_direction_180m"
	WindGusts10m             Metric = "wind_gusts_10m"
	Temperature80m           Metric = "temperature_80m"
	Temperature120m          Metric = "temperature_120m"
	Temperature180m          Metric = "temperature_180m"
	SoilTemperature0cm       Metric = "soil_temperature_0cm"
	SoilTemperature6cm       Metric = "soil_temperature_6cm"
	SoilTemperature18cm      Metric = "soil_temperature_18cm"
	SoilTemperature54cm      Metric = "soil_temperature_54cm"
	SoilMoisture0To1cm       Metric = "soil_moisture_0_to_1cm"
	SoilMoisture1To3cm       Metric = "soil_moisture_1_to_3cm"
	SoilMoisture3To9cm       Metric = "soil_moisture_3_to_9cm"
	SoilMoisture9To27cm      Metric = "soil_moisture_9_to_27cm"

	// Daily Metrics
	Temperature2mMax         Metric = "temperature_2m_max"
	Temperature2mMin         Metric = "temperature_2m_min"
	ApparentTemperatureMax   Metric = "apparent_temperature_max"
	ApparentTemperatureMin   Metric = "apparent_temperature_min"
	PrecipitationSum         Metric = "precipitation_sum"
	RainSum                  Metric = "rain_sum"
	ShowersSum               Metric = "showers_sum"
	SnowfallSum              Metric = "snowfall_sum"
	PrecipitationHours       Metric = "precipitation_hours"
	WeatherCodeDaily         Metric = "weather_code"
	Sunrise                  Metric = "sunrise"
	Sunset                   Metric = "sunset"
	WindSpeed10mMax          Metric = "wind_speed_10m_max"
	WindGusts10mMax          Metric = "wind_gusts_10m_max"
	WindDirection10mDominant Metric = "wind_direction_10m_dominant"
	ShortwaveRadiationSum    Metric = "shortwave_radiation_sum"

	// Current Metrics
	IsDay Metric = "is_day"

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

	// Marine Metrics (Hourly)
	WaveHeight                  Metric = "wave_height"
	WaveDirection               Metric = "wave_direction"
	WavePeriod                  Metric = "wave_period"
	WavePeakPeriod              Metric = "wave_peak_period"
	WindWaveHeight              Metric = "wind_wave_height"
	WindWaveDirection           Metric = "wind_wave_direction"
	WindWavePeriod              Metric = "wind_wave_period"
	WindWavePeakPeriod          Metric = "wind_wave_peak_period"
	SwellWaveHeight             Metric = "swell_wave_height"
	SwellWaveDirection          Metric = "swell_wave_direction"
	SwellWavePeriod             Metric = "swell_wave_period"
	SwellWavePeakPeriod         Metric = "swell_wave_peak_period"
	SecondarySwellWaveHeight    Metric = "secondary_swell_wave_height"
	SecondarySwellWaveDirection Metric = "secondary_swell_wave_direction"
	SecondarySwellWavePeriod    Metric = "secondary_swell_wave_period"
	TertiarySwellWaveHeight     Metric = "tertiary_swell_wave_height"
	TertiarySwellWaveDirection  Metric = "tertiary_swell_wave_direction"
	TertiarySwellWavePeriod     Metric = "tertiary_swell_wave_period"
	SeaLevelHeight              Metric = "sea_level_height"
	SeaSurfaceTemperature       Metric = "sea_surface_temperature"
	OceanCurrentVelocity        Metric = "ocean_current_velocity"
	OceanCurrentDirection       Metric = "ocean_current_direction"

	// Marine Metrics (Daily)
	WaveHeightMax              Metric = "wave_height_max"
	WaveDirectionDominant      Metric = "wave_direction_dominant"
	WavePeriodMax              Metric = "wave_period_max"
	WindWaveHeightMax          Metric = "wind_wave_height_max"
	WindWaveDirectionDominant  Metric = "wind_wave_direction_dominant"
	WindWavePeriodMax          Metric = "wind_wave_period_max"
	WindWavePeakPeriodMax      Metric = "wind_wave_peak_period_max"
	SwellWaveHeightMax         Metric = "swell_wave_height_max"
	SwellWaveDirectionDominant Metric = "swell_wave_direction_dominant"
	SwellWavePeriodMax         Metric = "swell_wave_period_max"
	SwellWavePeakPeriodMax     Metric = "swell_wave_peak_period_max"
)

```

#### **2. `openmeteo.go**`

I have added `GetMarine`, the `marineUrl` helper, and updated the structs to hold the new fields.

```go
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
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultHost          = "api.open-meteo.com"
	defaultScheme        = "https"
	seasonalHost         = "seasonal-api.open-meteo.com"
	marineHost           = "marine-api.open-meteo.com"
	forecastHistoryLimit = 7 * 24 * time.Hour

	// defaultUserAgent is the default User-Agent string sent with HTTP requests.
	defaultUserAgent = "OpenMeteoGo-Client"
)

// Client is used to interact with the Open-Meteo API.
type Client struct {
	// UserAgent is the string sent in the User-Agent header of the request.
	UserAgent string
	// HTTPClient allows for a custom http.Client to be used for requests.
	HTTPClient *http.Client
	apiKey     string
	scheme     string
	host       string
}

// Get fetches weather data based on the provided Options (Forecast and Archive).
func (c *Client) Get(o *Options) (*WeatherData, error) {
	req, err := http.NewRequest("GET", c.url(o), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server http error: %d", res.StatusCode)
	}

	var wd WeatherData
	if err = json.NewDecoder(res.Body).Decode(&wd); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &wd, nil
}

// GetSeasonal fetches seasonal weather forecast data.
func (c *Client) GetSeasonal(o *Options) (*WeatherData, error) {
	req, err := http.NewRequest("GET", c.seasonalUrl(o), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server http error: %d", res.StatusCode)
	}

	var wd WeatherData
	if err = json.NewDecoder(res.Body).Decode(&wd); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &wd, nil
}

// GetMarine fetches marine weather forecast data.
func (c *Client) GetMarine(o *Options) (*WeatherData, error) {
	req, err := http.NewRequest("GET", c.marineUrl(o), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server http error: %d", res.StatusCode)
	}

	var wd WeatherData
	if err = json.NewDecoder(res.Body).Decode(&wd); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &wd, nil
}

// NewClient creates a new Client with default settings.
func NewClient() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		UserAgent:  defaultUserAgent,
		scheme:     defaultScheme,
		host:       defaultHost,
	}
}

// NewClientWithKey creates a new Client configured with a commercial API key.
func NewClientWithKey(key string) *Client {
	return &Client{
		apiKey:     key,
		HTTPClient: http.DefaultClient,
		UserAgent:  defaultUserAgent,
		scheme:     defaultScheme,
		host:       defaultHost,
	}
}

func (c *Client) url(o *Options) string {
	host := c.host
	path := "/v1/forecast"

	// Determine if the request is for data older than the forecast API's history limit.
	isHistorical := !o.Start.IsZero() && time.Since(o.Start) > forecastHistoryLimit

	if isHistorical {
		host = "archive-" + host
		path = "/v1/archive"
	}

	if c.apiKey != "" {
		host = "customer-" + host
	}

	u := url.URL{
		Scheme: c.scheme,
		Host:   host,
		Path:   path,
	}

	q := u.Query()
	c.encodeCommonOptions(q, o)

	if o.PastDays > 0 {
		q.Set("past_days", fmt.Sprintf("%v", o.PastDays))
	}

	if o.ForcastDays > 0 {
		q.Set("forecast_days", fmt.Sprintf("%v", o.ForcastDays))
	}

	if o.HourlyMetrics != nil {
		if val := o.HourlyMetrics.encode(); val != "" {
			q.Set("hourly", val)
		}
	}

	if o.DailyMetrics != nil {
		if val := o.DailyMetrics.encode(); val != "" {
			q.Set("daily", val)
		}
	}

	if o.CurrentMetrics != nil {
		if val := o.CurrentMetrics.encode(); val != "" {
			q.Set("current", val)
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) seasonalUrl(o *Options) string {
	host := seasonalHost
	if c.apiKey != "" {
		host = "customer-" + host
	}

	u := url.URL{
		Scheme: c.scheme,
		Host:   host,
		Path:   "/v1/seasonal",
	}

	q := u.Query()
	c.encodeCommonOptions(q, o)

	if len(o.Models) > 0 {
		q.Set("models", strings.Join(o.Models, ","))
	}

	if o.HourlyMetrics != nil {
		if val := o.HourlyMetrics.encode(); val != "" {
			q.Set("hourly", val)
		}
	}
	if o.DailyMetrics != nil {
		if val := o.DailyMetrics.encode(); val != "" {
			q.Set("daily", val)
		}
	}
	if o.WeeklyMetrics != nil {
		if val := o.WeeklyMetrics.encode(); val != "" {
			q.Set("weekly", val)
		}
	}
	if o.MonthlyMetrics != nil {
		if val := o.MonthlyMetrics.encode(); val != "" {
			q.Set("monthly", val)
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) marineUrl(o *Options) string {
	host := marineHost
	if c.apiKey != "" {
		host = "customer-" + host
	}

	u := url.URL{
		Scheme: c.scheme,
		Host:   host,
		Path:   "/v1/marine",
	}

	q := u.Query()
	c.encodeCommonOptions(q, o)

	if len(o.Models) > 0 {
		q.Set("models", strings.Join(o.Models, ","))
	}

	if o.PastDays > 0 {
		q.Set("past_days", fmt.Sprintf("%v", o.PastDays))
	}

	if o.ForcastDays > 0 {
		q.Set("forecast_days", fmt.Sprintf("%v", o.ForcastDays))
	}

	if o.HourlyMetrics != nil {
		if val := o.HourlyMetrics.encode(); val != "" {
			q.Set("hourly", val)
		}
	}

	if o.DailyMetrics != nil {
		if val := o.DailyMetrics.encode(); val != "" {
			q.Set("daily", val)
		}
	}

	if o.CurrentMetrics != nil {
		if val := o.CurrentMetrics.encode(); val != "" {
			q.Set("current", val)
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// encodeCommonOptions encodes options shared between Forecast, Archive, Seasonal, and Marine APIs.
func (c *Client) encodeCommonOptions(q url.Values, o *Options) {
	if c.apiKey != "" {
		q.Set("apikey", c.apiKey)
	}

	q.Set("latitude", fmt.Sprintf("%v", o.Latitude))
	q.Set("longitude", fmt.Sprintf("%v", o.Longitude))

	if o.TemperatureUnit != "" {
		q.Set("temperature_unit", string(o.TemperatureUnit))
	}

	if o.WindspeedUnit != "" {
		q.Set("windspeed_unit", string(o.WindspeedUnit))
	}

	if o.PrecipitationUnit != "" {
		q.Set("precipitation_unit", string(o.PrecipitationUnit))
	}

	if o.Timezone.String() != "" {
		q.Set("timezone", o.Timezone.String())
	}

	if !o.Start.IsZero() {
		q.Set("start_date", o.Start.Format("2006-01-02"))
	}

	if !o.End.IsZero() {
		q.Set("end_date", o.End.Format("2006-01-02"))
	}
}

// WeatherData is the main struct that holds all the data returned from the API.
type WeatherData struct {
	Latitude             float64      `json:"latitude"`
	Longitude            float64      `json:"longitude"`
	GenerationtimeMs     float64      `json:"generationtime_ms"`
	UtcOffsetSeconds     int          `json:"utc_offset_seconds"`
	Timezone             string       `json:"timezone"`
	TimezoneAbbreviation string       `json:"timezone_abbreviation"`
	Elevation            float64      `json:"elevation"`
	CurrentUnits         CurrentUnits `json:"current_units"`
	Current              Current      `json:"current"`
	HourlyUnits          HourlyUnits  `json:"hourly_units"`
	Hourly               Hourly       `json:"hourly"`
	DailyUnits           DailyUnits   `json:"daily_units"`
	Daily                Daily        `json:"daily"`
	WeeklyUnits          WeeklyUnits  `json:"weekly_units"`
	Weekly               Weekly       `json:"weekly"`
	MonthlyUnits         MonthlyUnits `json:"monthly_units"`
	Monthly              Monthly      `json:"monthly"`
}

// CurrentUnits describes the units for the current weather data.
type CurrentUnits struct {
	Time                string `json:"time"`
	Interval            string `json:"interval"`
	Temperature2m       string `json:"temperature_2m"`
	RelativeHumidity2m  string `json:"relative_humidity_2m"`
	IsDay               string `json:"is_day"`
	ApparentTemperature string `json:"apparent_temperature"`
	Precipitation       string `json:"precipitation"`
	Rain                string `json:"rain"`
	Showers             string `json:"showers"`
	Snowfall            string `json:"snowfall"`
	WeatherCode         string `json:"weather_code"`
	CloudCover          string `json:"cloud_cover"`
	PressureMsl         string `json:"pressure_msl"`
	SurfacePressure     string `json:"surface_pressure"`
	WindSpeed10m        string `json:"wind_speed_10m"`
	WindDirection10m    string `json:"wind_direction_10m"`
	WindGusts10m        string `json:"wind_gusts_10m"`
}

// Current holds the current weather data.
type Current struct {
	Time                string  `json:"time"`
	Interval            int     `json:"interval"`
	Temperature2m       float64 `json:"temperature_2m"`
	RelativeHumidity2m  int     `json:"relative_humidity_2m"`
	IsDay               int     `json:"is_day"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	Precipitation       float64 `json:"precipitation"`
	Rain                float64 `json:"rain"`
	Showers             float64 `json:"showers"`
	Snowfall            float64 `json:"snowfall"`
	WeatherCode         int     `json:"weather_code"`
	CloudCover          int     `json:"cloud_cover"`
	PressureMsl         float64 `json:"pressure_msl"`
	SurfacePressure     float64 `json:"surface_pressure"`
	WindSpeed10m        float64 `json:"wind_speed_10m"`
	WindDirection10m    int     `json:"wind_direction_10m"`
	WindGusts10m        float64 `json:"wind_gusts_10m"`
}

// HourlyUnits describes the units for the hourly forecast data.
type HourlyUnits struct {
	Time                        string `json:"time"`
	Temperature2m               string `json:"temperature_2m"`
	RelativeHumidity2m          string `json:"relative_humidity_2m"`
	DewPoint2m                  string `json:"dew_point_2m"`
	ApparentTemperature         string `json:"apparent_temperature"`
	PrecipitationProbability    string `json:"precipitation_probability"`
	Precipitation               string `json:"precipitation"`
	Rain                        string `json:"rain"`
	Showers                     string `json:"showers"`
	Snowfall                    string `json:"snowfall"`
	SnowDepth                   string `json:"snow_depth"`
	WeatherCode                 string `json:"weather_code"`
	PressureMsl                 string `json:"pressure_msl"`
	SurfacePressure             string `json:"surface_pressure"`
	CloudCover                  string `json:"cloud_cover"`
	CloudCoverLow               string `json:"cloud_cover_low"`
	CloudCoverMid               string `json:"cloud_cover_mid"`
	CloudCoverHigh              string `json:"cloud_cover_high"`
	Evapotranspiration          string `json:"evapotranspiration"`
	Visibility                  string `json:"visibility"`
	Et0FaoEvapotranspiration    string `json:"et0_fao_evapotranspiration"`
	VapourPressureDeficit       string `json:"vapour_pressure_deficit"`
	WindSpeed10m                string `json:"wind_speed_10m"`
	WindSpeed80m                string `json:"wind_speed_80m"`
	WindSpeed120m               string `json:"wind_speed_120m"`
	WindSpeed180m               string `json:"wind_speed_180m"`
	WindDirection10m            string `json:"wind_direction_10m"`
	WindDirection80m            string `json:"wind_direction_80m"`
	WindDirection120m           string `json:"wind_direction_120m"`
	WindDirection180m           string `json:"wind_direction_180m"`
	WindGusts10m                string `json:"wind_gusts_10m"`
	Temperature80m              string `json:"temperature_80m"`
	Temperature120m             string `json:"temperature_120m"`
	Temperature180m             string `json:"temperature_180m"`
	SoilTemperature0cm          string `json:"soil_temperature_0cm"`
	SoilTemperature6cm          string `json:"soil_temperature_6cm"`
	SoilTemperature18cm         string `json:"soil_temperature_18cm"`
	SoilTemperature54cm         string `json:"soil_temperature_54cm"`
	SoilMoisture0To1cm          string `json:"soil_moisture_0_to_1cm"`
	SoilMoisture1To3cm          string `json:"soil_moisture_1_to_3cm"`
	SoilMoisture3To9cm          string `json:"soil_moisture_3_to_9cm"`
	SoilMoisture9To27cm         string `json:"soil_moisture_9_to_27cm"`
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
	SeaLevelHeight              string `json:"sea_level_height"`
	SeaSurfaceTemperature       string `json:"sea_surface_temperature"`
	OceanCurrentVelocity        string `json:"ocean_current_velocity"`
	OceanCurrentDirection       string `json:"ocean_current_direction"`
}

// Hourly holds slices for each hourly forecast metric.
type Hourly struct {
	Time                        []string  `json:"time"`
	Temperature2m               []float64 `json:"temperature_2m"`
	RelativeHumidity2m          []int     `json:"relative_humidity_2m"`
	DewPoint2m                  []float64 `json:"dew_point_2m"`
	ApparentTemperature         []float64 `json:"apparent_temperature"`
	PrecipitationProbability    []int     `json:"precipitation_probability"`
	Precipitation               []float64 `json:"precipitation"`
	Rain                        []float64 `json:"rain"`
	Showers                     []float64 `json:"showers"`
	Snowfall                    []float64 `json:"snowfall"`
	SnowDepth                   []float64 `json:"snow_depth"`
	WeatherCode                 []int     `json:"weather_code"`
	PressureMsl                 []float64 `json:"pressure_msl"`
	SurfacePressure             []float64 `json:"surface_pressure"`
	CloudCover                  []int     `json:"cloud_cover"`
	CloudCoverLow               []int     `json:"cloud_cover_low"`
	CloudCoverMid               []int     `json:"cloud_cover_mid"`
	CloudCoverHigh              []int     `json:"cloud_cover_high"`
	Evapotranspiration          []float64 `json:"evapotranspiration"`
	Visibility                  []float64 `json:"visibility"`
	Et0FaoEvapotranspiration    []float64 `json:"et0_fao_evapotranspiration"`
	VapourPressureDeficit       []float64 `json:"vapour_pressure_deficit"`
	WindSpeed10m                []float64 `json:"wind_speed_10m"`
	WindSpeed80m                []float64 `json:"wind_speed_80m"`
	WindSpeed120m               []float64 `json:"wind_speed_120m"`
	WindSpeed180m               []float64 `json:"wind_speed_180m"`
	WindDirection10m            []int     `json:"wind_direction_10m"`
	WindDirection80m            []int     `json:"wind_direction_80m"`
	WindDirection120m           []int     `json:"wind_direction_120m"`
	WindDirection180m           []int     `json:"wind_direction_180m"`
	WindGusts10m                []float64 `json:"wind_gusts_10m"`
	Temperature80m              []float64 `json:"temperature_80m"`
	Temperature120m             []float64 `json:"temperature_120m"`
	Temperature180m             []float64 `json:"temperature_180m"`
	SoilTemperature0cm          []float64 `json:"soil_temperature_0cm"`
	SoilTemperature6cm          []float64 `json:"soil_temperature_6cm"`
	SoilTemperature18cm         []float64 `json:"soil_temperature_18cm"`
	SoilTemperature54cm         []float64 `json:"soil_temperature_54cm"`
	SoilMoisture0To1cm          []float64 `json:"soil_moisture_0_to_1cm"`
	SoilMoisture1To3cm          []float64 `json:"soil_moisture_1_to_3cm"`
	SoilMoisture3To9cm          []float64 `json:"soil_moisture_3_to_9cm"`
	SoilMoisture9To27cm         []float64 `json:"soil_moisture_9_to_27cm"`
	WaveHeight                  []float64 `json:"wave_height"`
	WaveDirection               []float64 `json:"wave_direction"`
	WavePeriod                  []float64 `json:"wave_period"`
	WavePeakPeriod              []float64 `json:"wave_peak_period"`
	WindWaveHeight              []float64 `json:"wind_wave_height"`
	WindWaveDirection           []float64 `json:"wind_wave_direction"`
	WindWavePeriod              []float64 `json:"wind_wave_period"`
	WindWavePeakPeriod          []float64 `json:"wind_wave_peak_period"`
	SwellWaveHeight             []float64 `json:"swell_wave_height"`
	SwellWaveDirection          []float64 `json:"swell_wave_direction"`
	SwellWavePeriod             []float64 `json:"swell_wave_period"`
	SwellWavePeakPeriod         []float64 `json:"swell_wave_peak_period"`
	SecondarySwellWaveHeight    []float64 `json:"secondary_swell_wave_height"`
	SecondarySwellWaveDirection []float64 `json:"secondary_swell_wave_direction"`
	SecondarySwellWavePeriod    []float64 `json:"secondary_swell_wave_period"`
	TertiarySwellWaveHeight     []float64 `json:"tertiary_swell_wave_height"`
	TertiarySwellWaveDirection  []float64 `json:"tertiary_swell_wave_direction"`
	TertiarySwellWavePeriod     []float64 `json:"tertiary_swell_wave_period"`
	SeaLevelHeight              []float64 `json:"sea_level_height"`
	SeaSurfaceTemperature       []float64 `json:"sea_surface_temperature"`
	OceanCurrentVelocity        []float64 `json:"ocean_current_velocity"`
	OceanCurrentDirection       []float64 `json:"ocean_current_direction"`
}

// DailyUnits describes the units for the daily forecast data.
type DailyUnits struct {
	Time                       string `json:"time"`
	WeatherCode                string `json:"weather_code"`
	Temperature2mMax           string `json:"temperature_2m_max"`
	Temperature2mMin           string `json:"temperature_2m_min"`
	ApparentTemperatureMax     string `json:"apparent_temperature_max"`
	ApparentTemperatureMin     string `json:"apparent_temperature_min"`
	Sunrise                    string `json:"sunrise"`
	Sunset                     string `json:"sunset"`
	PrecipitationSum           string `json:"precipitation_sum"`
	RainSum                    string `json:"rain_sum"`
	ShowersSum                 string `json:"showers_sum"`
	SnowfallSum                string `json:"snowfall_sum"`
	PrecipitationHours         string `json:"precipitation_hours"`
	WindSpeed10mMax            string `json:"wind_speed_10m_max"`
	WindGusts10mMax            string `json:"wind_gusts_10m_max"`
	WindDirection10mDominant   string `json:"wind_direction_10m_dominant"`
	ShortwaveRadiationSum      string `json:"shortwave_radiation_sum"`
	Et0FaoEvapotranspiration   string `json:"et0_fao_evapotranspiration"`
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

// Daily holds slices for each daily forecast metric.
type Daily struct {
	Time                       []string  `json:"time"`
	WeatherCode                []int     `json:"weather_code"`
	Temperature2mMax           []float64 `json:"temperature_2m_max"`
	Temperature2mMin           []float64 `json:"temperature_2m_min"`
	ApparentTemperatureMax     []float64 `json:"apparent_temperature_max"`
	ApparentTemperatureMin     []float64 `json:"apparent_temperature_min"`
	Sunrise                    []string  `json:"sunrise"`
	Sunset                     []string  `json:"sunset"`
	PrecipitationSum           []float64 `json:"precipitation_sum"`
	RainSum                    []float64 `json:"rain_sum"`
	ShowersSum                 []float64 `json:"showers_sum"`
	SnowfallSum                []float64 `json:"snowfall_sum"`
	PrecipitationHours         []float64 `json:"precipitation_hours"`
	WindSpeed10mMax            []float64 `json:"wind_speed_10m_max"`
	WindGusts10mMax            []float64 `json:"wind_gusts_10m_max"`
	WindDirection10mDominant   []int     `json:"wind_direction_10m_dominant"`
	ShortwaveRadiationSum      []float64 `json:"shortwave_radiation_sum"`
	Et0FaoEvapotranspiration   []float64 `json:"et0_fao_evapotranspiration"`
	WaveHeightMax              []float64 `json:"wave_height_max"`
	WaveDirectionDominant      []float64 `json:"wave_direction_dominant"`
	WavePeriodMax              []float64 `json:"wave_period_max"`
	WindWaveHeightMax          []float64 `json:"wind_wave_height_max"`
	WindWaveDirectionDominant  []float64 `json:"wind_wave_direction_dominant"`
	WindWavePeriodMax          []float64 `json:"wind_wave_period_max"`
	WindWavePeakPeriodMax      []float64 `json:"wind_wave_peak_period_max"`
	SwellWaveHeightMax         []float64 `json:"swell_wave_height_max"`
	SwellWaveDirectionDominant []float64 `json:"swell_wave_direction_dominant"`
	SwellWavePeriodMax         []float64 `json:"swell_wave_period_max"`
	SwellWavePeakPeriodMax     []float64 `json:"swell_wave_peak_period_max"`
}

// WeeklyUnits describes the units for the weekly seasonal forecast data.
type WeeklyUnits struct {
	Time                    string `json:"time"`
	Temperature2mMean       string `json:"temperature_2m_mean"`
	Temperature2mAnomaly    string `json:"temperature_2m_anomaly"`
	PrecipitationMean       string `json:"precipitation_mean"`
	PrecipitationAnomaly    string `json:"precipitation_anomaly"`
	PressureMslMean         string `json:"pressure_msl_mean"`
	PressureMslAnomaly      string `json:"pressure_msl_anomaly"`
	SoilMoisture0To10cmMean string `json:"soil_moisture_0_to_10cm_mean"`
}

// Weekly holds slices for each weekly seasonal forecast metric.
type Weekly struct {
	Time                    []string  `json:"time"`
	Temperature2mMean       []float64 `json:"temperature_2m_mean"`
	Temperature2mAnomaly    []float64 `json:"temperature_2m_anomaly"`
	PrecipitationMean       []float64 `json:"precipitation_mean"`
	PrecipitationAnomaly    []float64 `json:"precipitation_anomaly"`
	PressureMslMean         []float64 `json:"pressure_msl_mean"`
	PressureMslAnomaly      []float64 `json:"pressure_msl_anomaly"`
	SoilMoisture0To10cmMean []float64 `json:"soil_moisture_0_to_10cm_mean"`
}

// MonthlyUnits describes the units for the monthly seasonal forecast data.
type MonthlyUnits struct {
	Time                    string `json:"time"`
	Temperature2mMean       string `json:"temperature_2m_mean"`
	Temperature2mAnomaly    string `json:"temperature_2m_anomaly"`
	PrecipitationMean       string `json:"precipitation_mean"`
	PrecipitationAnomaly    string `json:"precipitation_anomaly"`
	PressureMslMean         string `json:"pressure_msl_mean"`
	PressureMslAnomaly      string `json:"pressure_msl_anomaly"`
	SoilMoisture0To10cmMean string `json:"soil_moisture_0_to_10cm_mean"`
}

// Monthly holds slices for each monthly seasonal forecast metric.
type Monthly struct {
	Time                    []string  `json:"time"`
	Temperature2mMean       []float64 `json:"temperature_2m_mean"`
	Temperature2mAnomaly    []float64 `json:"temperature_2m_anomaly"`
	PrecipitationMean       []float64 `json:"precipitation_mean"`
	PrecipitationAnomaly    []float64 `json:"precipitation_anomaly"`
	PressureMslMean         []float64 `json:"pressure_msl_mean"`
	PressureMslAnomaly      []float64 `json:"pressure_msl_anomaly"`
	SoilMoisture0To10cmMean []float64 `json:"soil_moisture_0_to_10cm_mean"`
}

```