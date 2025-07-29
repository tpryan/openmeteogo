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
	"time"
)

const (
	defaultHost          = "api.open-meteo.com"
	defaultScheme        = "https"
	forecastHistoryLimit = 90 * 24 * time.Hour

	// DefaultUserAgent is the default User-Agent string sent with HTTP requests.
	DefaultUserAgent = "OpenMeteoGo-Client"
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

// Get fetches weather data based on the provided Options.
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

// NewClient creates a new Client with default settings.
func NewClient() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		UserAgent:  DefaultUserAgent,
		scheme:     defaultScheme,
		host:       defaultHost,
	}
}

// NewClientWithKey creates a new Client configured with a commercial API key.
func NewClientWithKey(key string) *Client {
	return &Client{
		apiKey:     key,
		HTTPClient: http.DefaultClient,
		UserAgent:  DefaultUserAgent,
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

	if c.apiKey != "" {
		q.Set("apikey", c.apiKey)
	}

	q.Set("latitude", fmt.Sprintf("%v", o.Latitude))
	q.Set("longitude", fmt.Sprintf("%v", o.Longitude))

	if o.TemperatureUnit > 0 {
		q.Set("temperature_unit", o.TemperatureUnit.String())
	}

	if o.WindspeedUnit > 0 {
		q.Set("windspeed_unit", o.WindspeedUnit.String())
	}

	if o.PrecipitationUnit > 0 {
		q.Set("precipitation_unit", o.PrecipitationUnit.String())
	}

	if o.Timezone.String() != "" {
		q.Set("timezone", o.Timezone.String())
	}

	if o.PastDays > 0 {
		q.Set("past_days", fmt.Sprintf("%v", o.PastDays))
	}

	if !o.Start.IsZero() {
		q.Set("start_date", o.Start.Format("2006-01-02"))
	}

	if !o.End.IsZero() {
		q.Set("end_date", o.End.Format("2006-01-02"))
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

	if o.ForcastDays > 0 {
		q.Set("forecast_days", fmt.Sprintf("%v", o.ForcastDays))
	}

	if o.CurrentMetrics != nil {
		if val := o.CurrentMetrics.encode(); val != "" {
			q.Set("current", val)
		}
	}

	u.RawQuery = q.Encode()

	return u.String()
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

// Current holds the current weather data values.
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
	Time                     string `json:"time"`
	Temperature2m            string `json:"temperature_2m"`
	RelativeHumidity2m       string `json:"relative_humidity_2m"`
	DewPoint2m               string `json:"dew_point_2m"`
	ApparentTemperature      string `json:"apparent_temperature"`
	PrecipitationProbability string `json:"precipitation_probability"`
	Precipitation            string `json:"precipitation"`
	Rain                     string `json:"rain"`
	Showers                  string `json:"showers"`
	Snowfall                 string `json:"snowfall"`
	SnowDepth                string `json:"snow_depth"`
	WeatherCode              string `json:"weather_code"`
	PressureMsl              string `json:"pressure_msl"`
	SurfacePressure          string `json:"surface_pressure"`
	CloudCover               string `json:"cloud_cover"`
	CloudCoverLow            string `json:"cloud_cover_low"`
	CloudCoverMid            string `json:"cloud_cover_mid"`
	CloudCoverHigh           string `json:"cloud_cover_high"`
	Evapotranspiration       string `json:"evapotranspiration"`
	Visibility               string `json:"visibility"`
	Et0FaoEvapotranspiration string `json:"et0_fao_evapotranspiration"`
	VapourPressureDeficit    string `json:"vapour_pressure_deficit"`
	WindSpeed10m             string `json:"wind_speed_10m"`
	WindSpeed80m             string `json:"wind_speed_80m"`
	WindSpeed120m            string `json:"wind_speed_120m"`
	WindSpeed180m            string `json:"wind_speed_180m"`
	WindDirection10m         string `json:"wind_direction_10m"`
	WindDirection80m         string `json:"wind_direction_80m"`
	WindDirection120m        string `json:"wind_direction_120m"`
	WindDirection180m        string `json:"wind_direction_180m"`
	WindGusts10m             string `json:"wind_gusts_10m"`
	Temperature80m           string `json:"temperature_80m"`
	Temperature120m          string `json:"temperature_120m"`
	Temperature180m          string `json:"temperature_180m"`
	SoilTemperature0cm       string `json:"soil_temperature_0cm"`
	SoilTemperature6cm       string `json:"soil_temperature_6cm"`
	SoilTemperature18cm      string `json:"soil_temperature_18cm"`
	SoilTemperature54cm      string `json:"soil_temperature_54cm"`
	SoilMoisture0To1cm       string `json:"soil_moisture_0_to_1cm"`
	SoilMoisture1To3cm       string `json:"soil_moisture_1_to_3cm"`
	SoilMoisture9To27cm      string `json:"soil_moisture_9_to_27cm"`
	SoilMoisture3To9cm       string `json:"soil_moisture_3_to_9cm"`
}

// Hourly holds slices for each hourly forecast metric.
type Hourly struct {
	Time                     []string  `json:"time"`
	Temperature2m            []float64 `json:"temperature_2m"`
	RelativeHumidity2m       []int     `json:"relative_humidity_2m"`
	DewPoint2m               []float64 `json:"dew_point_2m"`
	ApparentTemperature      []float64 `json:"apparent_temperature"`
	PrecipitationProbability []int     `json:"precipitation_probability"`
	Precipitation            []float64 `json:"precipitation"`
	Rain                     []float64 `json:"rain"`
	Showers                  []float64 `json:"showers"`
	Snowfall                 []float64 `json:"snowfall"`
	SnowDepth                []float64 `json:"snow_depth"`
	WeatherCode              []int     `json:"weather_code"`
	PressureMsl              []float64 `json:"pressure_msl"`
	SurfacePressure          []float64 `json:"surface_pressure"`
	CloudCover               []int     `json:"cloud_cover"`
	CloudCoverLow            []int     `json:"cloud_cover_low"`
	CloudCoverMid            []int     `json:"cloud_cover_mid"`
	CloudCoverHigh           []int     `json:"cloud_cover_high"`
	Evapotranspiration       []float64 `json:"evapotranspiration"`
	Visibility               []float64 `json:"visibility"`
	Et0FaoEvapotranspiration []float64 `json:"et0_fao_evapotranspiration"`
	VapourPressureDeficit    []float64 `json:"vapour_pressure_deficit"`
	WindSpeed10m             []float64 `json:"wind_speed_10m"`
	WindSpeed80m             []float64 `json:"wind_speed_80m"`
	WindSpeed120m            []float64 `json:"wind_speed_120m"`
	WindSpeed180m            []float64 `json:"wind_speed_180m"`
	WindDirection10m         []int     `json:"wind_direction_10m"`
	WindDirection80m         []int     `json:"wind_direction_80m"`
	WindDirection120m        []int     `json:"wind_direction_120m"`
	WindDirection180m        []int     `json:"wind_direction_180m"`
	WindGusts10m             []float64 `json:"wind_gusts_10m"`
	Temperature80m           []float64 `json:"temperature_80m"`
	Temperature120m          []float64 `json:"temperature_120m"`
	Temperature180m          []float64 `json:"temperature_180m"`
	SoilTemperature0cm       []float64 `json:"soil_temperature_0cm"`
	SoilTemperature6cm       []float64 `json:"soil_temperature_6cm"`
	SoilTemperature18cm      []float64 `json:"soil_temperature_18cm"`
	SoilTemperature54cm      []float64 `json:"soil_temperature_54cm"`
	SoilMoisture0To1cm       []float64 `json:"soil_moisture_0_to_1cm"`
	SoilMoisture1To3cm       []float64 `json:"soil_moisture_1_to_3cm"`
	SoilMoisture9To27cm      []float64 `json:"soil_moisture_9_to_27cm"`
	SoilMoisture3To9cm       []float64 `json:"soil_moisture_3_to_9cm"`
}

// DailyUnits describes the units for the daily forecast data.
type DailyUnits struct {
	Time                        string `json:"time"`
	WeatherCode                 string `json:"weather_code"`
	Temperature2mMax            string `json:"temperature_2m_max"`
	Temperature2mMin            string `json:"temperature_2m_min"`
	ApparentTemperatureMax      string `json:"apparent_temperature_max"`
	ApparentTemperatureMin      string `json:"apparent_temperature_min"`
	Sunrise                     string `json:"sunrise"`
	Sunset                      string `json:"sunset"`
	SunshineDuration            string `json:"sunshine_duration"`
	DaylightDuration            string `json:"daylight_duration"`
	UvIndexMax                  string `json:"uv_index_max"`
	UvIndexClearSkyMax          string `json:"uv_index_clear_sky_max"`
	RainSum                     string `json:"rain_sum"`
	ShowersSum                  string `json:"showers_sum"`
	SnowfallSum                 string `json:"snowfall_sum"`
	PrecipitationSum            string `json:"precipitation_sum"`
	PrecipitationHours          string `json:"precipitation_hours"`
	PrecipitationProbabilityMax string `json:"precipitation_probability_max"`
	WindSpeed10mMax             string `json:"wind_speed_10m_max"`
	WindGusts10mMax             string `json:"wind_gusts_10m_max"`
	WindDirection10mDominant    string `json:"wind_direction_10m_dominant"`
	ShortwaveRadiationSum       string `json:"shortwave_radiation_sum"`
	Et0FaoEvapotranspiration    string `json:"et0_fao_evapotranspiration"`
}

// Daily holds slices for each daily forecast metric.
type Daily struct {
	Time                        []string  `json:"time"`
	WeatherCode                 []int     `json:"weather_code"`
	Temperature2mMax            []float64 `json:"temperature_2m_max"`
	Temperature2mMin            []float64 `json:"temperature_2m_min"`
	ApparentTemperatureMax      []float64 `json:"apparent_temperature_max"`
	ApparentTemperatureMin      []float64 `json:"apparent_temperature_min"`
	Sunrise                     []string  `json:"sunrise"`
	Sunset                      []string  `json:"sunset"`
	SunshineDuration            []float64 `json:"sunshine_duration"`
	DaylightDuration            []float64 `json:"daylight_duration"`
	UvIndexMax                  []float64 `json:"uv_index_max"`
	UvIndexClearSkyMax          []float64 `json:"uv_index_clear_sky_max"`
	RainSum                     []float64 `json:"rain_sum"`
	ShowersSum                  []float64 `json:"showers_sum"`
	SnowfallSum                 []float64 `json:"snowfall_sum"`
	PrecipitationSum            []float64 `json:"precipitation_sum"`
	PrecipitationHours          []float64 `json:"precipitation_hours"`
	PrecipitationProbabilityMax []int     `json:"precipitation_probability_max"`
	WindSpeed10mMax             []float64 `json:"wind_speed_10m_max"`
	WindGusts10mMax             []float64 `json:"wind_gusts_10m_max"`
	WindDirection10mDominant    []int     `json:"wind_direction_10m_dominant"`
	ShortwaveRadiationSum       []float64 `json:"shortwave_radiation_sum"`
	Et0FaoEvapotranspiration    []float64 `json:"et0_fao_evapotranspiration"`
}
