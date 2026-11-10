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
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHost         = "api.open-meteo.com"
	defaultScheme       = "https"
	defaultSeasonalHost = "seasonal-api.open-meteo.com"
	defaultMarineHost   = "marine-api.open-meteo.com"
	// seasonalPath is the path of the seasonal endpoint.
	seasonalPath         = "/v1/seasonal"
	forecastHistoryLimit = 7 * 24 * time.Hour

	// defaultUserAgent is the default User-Agent string sent with HTTP requests.
	defaultUserAgent = "OpenMeteoGo-Client"

	// defaultTimeout bounds a single request. http.DefaultClient has no
	// timeout, so without this a stalled connection hangs the caller forever.
	defaultTimeout = 30 * time.Second

	// errorBodyLimit caps how much of an error response body is read. The API
	// sends a short JSON object, but a proxy in the path can return anything.
	errorBodyLimit = 64 << 10
)

// APIError reports a non-200 response from the Open-Meteo API.
//
// The service describes failures in the response body as
// {"error": true, "reason": "..."}, and Reason carries that text so callers
// can see why a request was rejected instead of only its status code.
type APIError struct {
	// StatusCode is the HTTP status returned by the API.
	StatusCode int
	// Reason is the explanation taken from the response body, empty when the
	// body carried none.
	Reason string
	// RetryAfter is the delay requested by a Retry-After header, which the API
	// sends when it rate limits a client. Zero when the header was absent.
	RetryAfter time.Duration
}

func (e *APIError) Error() string {
	switch {
	case e.Reason != "" && e.RetryAfter > 0:
		return fmt.Sprintf("open-meteo: http %d: %s (retry after %s)", e.StatusCode, e.Reason, e.RetryAfter)
	case e.Reason != "":
		return fmt.Sprintf("open-meteo: http %d: %s", e.StatusCode, e.Reason)
	default:
		return fmt.Sprintf("open-meteo: http %d", e.StatusCode)
	}
}

// RateLimited reports whether the request was rejected for rate limiting.
// Pair it with RetryAfter to decide how long to wait before retrying.
func (e *APIError) RateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// newAPIError builds an APIError from a failed response, reading the body for
// the reason the API gives for the rejection.
func newAPIError(res *http.Response) *APIError {
	e := &APIError{StatusCode: res.StatusCode}

	data, err := io.ReadAll(io.LimitReader(res.Body, errorBodyLimit))
	if err == nil && len(data) > 0 {
		// Decode only the reason. The "error" flag adds nothing a non-200
		// status has not already said, and it is not always a JSON boolean.
		var body struct {
			Reason string `json:"reason"`
		}
		if json.Unmarshal(data, &body) == nil {
			e.Reason = body.Reason
		} else if txt := strings.TrimSpace(string(data)); txt != "" {
			// Not JSON at all: most likely a gateway error page.
			e.Reason = truncate(txt, 200)
		}
	}

	if v := res.Header.Get("Retry-After"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs >= 0 {
			e.RetryAfter = time.Duration(secs) * time.Second
		} else if t, err := http.ParseTime(v); err == nil {
			if d := time.Until(t); d > 0 {
				e.RetryAfter = d
			}
		}
	}

	return e
}

// truncate shortens s to at most n runes, never splitting one.
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}

	return string(r[:n]) + "…"
}

// Client is used to interact with the Open-Meteo API.
type Client struct {
	// UserAgent is the string sent in the User-Agent header of the request.
	UserAgent string
	// HTTPClient allows for a custom http.Client to be used for requests.
	HTTPClient   *http.Client
	apiKey       string
	scheme       string
	host         string
	seasonalHost string
	marineHost   string
}

// Get fetches weather data based on the provided Options.
//
// The context bounds the whole request: cancelling it aborts an in-flight
// call. A non-200 response is returned as an *APIError carrying the reason
// the API gave, so use errors.As to inspect it.
func (c *Client) Get(ctx context.Context, o *Options) (*WeatherData, error) {
	if o == nil {
		return nil, errors.New("openmeteogo: nil Options")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(o), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", cmp.Or(c.UserAgent, defaultUserAgent))

	// The key travels in a header rather than the query string. Open-Meteo
	// accepts both, but a transport failure returns a *url.Error holding the
	// full URL, which would put the key straight into the caller's logs.
	if c.apiKey != "" {
		req.Header.Set("X-Api-Key", c.apiKey)
	}

	res, err := c.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, newAPIError(res)
	}

	var wd WeatherData
	if err = json.NewDecoder(res.Body).Decode(&wd); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &wd, nil
}

// httpClient returns the client to send with, falling back to a timeout-bearing
// default so a zero-value Client works instead of panicking.
func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}

	return newHTTPClient()
}

// newHTTPClient builds the default transport. It is deliberately not
// http.DefaultClient: that one has no timeout, and it is process-global, so
// setting a timeout on it would reach into unrelated code.
func newHTTPClient() *http.Client {
	return &http.Client{Timeout: defaultTimeout}
}

// NewClient creates a new Client with default settings.
func NewClient() *Client {
	return &Client{
		HTTPClient:   newHTTPClient(),
		UserAgent:    defaultUserAgent,
		scheme:       defaultScheme,
		host:         defaultHost,
		seasonalHost: defaultSeasonalHost,
		marineHost:   defaultMarineHost,
	}
}

// NewClientWithKey creates a new Client configured with a commercial API key.
func NewClientWithKey(key string) *Client {
	c := NewClient()
	c.apiKey = key

	return c
}

func (c *Client) url(o *Options) string {
	// Fall back to the defaults so a Client built as a struct literal, rather
	// than through NewClient, still produces a usable URL.
	host := cmp.Or(c.host, defaultHost)
	path := "/v1/forecast"

	// Determine if the request is for seasonal data. Models is deliberately
	// not part of this test: the forecast and marine endpoints accept a
	// models parameter too, so keying off it would silently reroute an
	// ordinary forecast request to the seasonal API.
	isSeasonal := o.Seasonal || len(o.WeeklyMetrics) > 0 || len(o.MonthlyMetrics) > 0

	// Determine if the request is for marine data.
	isMarine := o.Marine

	// Determine if the request is for data older than the forecast API's
	// history limit. The whole range has to be outside the window: if End
	// reaches back into it, the forecast endpoint still serves the request,
	// and rewriting to the archive would silently drop the newer half.
	isHistorical := false
	if !o.Start.IsZero() && time.Since(o.Start) > forecastHistoryLimit {
		isHistorical = o.End.IsZero() || time.Since(o.End) > forecastHistoryLimit
	}

	// The marine and seasonal endpoints each serve their own history, so they
	// win over the archive rewrite rather than combining with it.
	switch {
	case isMarine:
		host = cmp.Or(c.marineHost, defaultMarineHost)
		path = marinePath
	case isSeasonal:
		host = cmp.Or(c.seasonalHost, defaultSeasonalHost)
		path = seasonalPath
	case isHistorical:
		host = "archive-" + host
		path = "/v1/archive"
	}

	if c.apiKey != "" {
		host = "customer-" + host
	}

	u := url.URL{
		Scheme: cmp.Or(c.scheme, defaultScheme),
		Host:   host,
		Path:   path,
	}

	q := u.Query()

	// Use common options encoding
	c.encodeCommonOptions(q, o)

	// Every endpoint accepts a model selection, so this is not gated on the
	// request being seasonal.
	if len(o.Models) > 0 {
		q.Set("models", strings.Join(o.Models, ","))
	}

	if isSeasonal {
		if val := o.WeeklyMetrics.encode(); val != "" {
			q.Set("weekly", val)
		}
		if val := o.MonthlyMetrics.encode(); val != "" {
			q.Set("monthly", val)
		}
	}

	if val := o.HourlyMetrics.encode(); val != "" {
		q.Set("hourly", val)
	}

	if val := o.DailyMetrics.encode(); val != "" {
		q.Set("daily", val)
	}

	if o.ForecastDays > 0 {
		q.Set("forecast_days", strconv.Itoa(o.ForecastDays))
	}

	if val := o.CurrentMetrics.encode(); val != "" {
		q.Set("current", val)
	}

	u.RawQuery = q.Encode()

	return u.String()
}

// encodeCommonOptions encodes options shared between Forecast, Archive, and Seasonal APIs.
func (c *Client) encodeCommonOptions(q url.Values, o *Options) {
	q.Set("latitude", formatCoord(o.Latitude))
	q.Set("longitude", formatCoord(o.Longitude))

	if o.TemperatureUnit != "" {
		q.Set("temperature_unit", string(o.TemperatureUnit))
	}

	if o.WindspeedUnit != "" {
		q.Set("windspeed_unit", string(o.WindspeedUnit))
	}

	if o.PrecipitationUnit != "" {
		q.Set("precipitation_unit", string(o.PrecipitationUnit))
	}

	if o.Timezone != nil {
		if name := o.Timezone.String(); name != "" {
			q.Set("timezone", name)
		}
	}

	if !o.Start.IsZero() {
		q.Set("start_date", o.Start.Format("2006-01-02"))
	}

	if !o.End.IsZero() {
		q.Set("end_date", o.End.Format("2006-01-02"))
	}

	if o.PastDays > 0 {
		q.Set("past_days", strconv.Itoa(o.PastDays))
	}
}

// formatCoord renders a coordinate in plain decimal. The %v verb switches to
// scientific notation for small magnitudes, which the API rejects.
func formatCoord(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
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

	// CurrentMarineUnits carries the marine units, populated only by marine
	// requests. Its fields are promoted, so they read as units.WaveHeight.
	CurrentMarineUnits
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

	// CurrentMarine carries the marine values, populated only by marine
	// requests. Its fields are promoted, so they read as current.WaveHeight,
	// and each is a *float64 so a JSON null stays distinguishable from a
	// genuine zero.
	CurrentMarine
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

	// HourlyMarineUnits carries the marine units, populated only by marine
	// requests. Its fields are promoted, so they read as units.WaveHeight.
	HourlyMarineUnits
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

	// HourlyMarine carries the marine series, populated only by marine
	// requests. Its fields are promoted, so they read as hourly.WaveHeight,
	// and each is a []*float64 so a JSON null stays distinguishable from a
	// genuine zero.
	HourlyMarine
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

	// DailyMarineUnits carries the marine units, populated only by marine
	// requests. Its fields are promoted, so they read as units.WaveHeightMax.
	DailyMarineUnits
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

	// DailyMarine carries the marine series, populated only by marine
	// requests. Its fields are promoted, so they read as daily.WaveHeightMax,
	// and each is a []*float64 so a JSON null stays distinguishable from a
	// genuine zero.
	DailyMarine
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
