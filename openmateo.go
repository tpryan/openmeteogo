package openmateo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const baseHost = "api.open-meteo.com"

const DefaultUserAgent = "OpenMeteoGo-Client"

const TempUnitCelsius = "celsius"
const TempUnitFahrenheit = "fahrenheit"

const WindSpeedUnitKMH = "kmh"
const WindSpeedUnitMS = "ms"
const WindSpeedUnitMPH = "mph"
const WindspeedUnitKN = "kn"

const PrecipitationUnitMM = "mm"
const PrecipitationUnitIN = "in"

type Client struct {
	apiKey     string
	UserAgent  string
	HTTPClient *http.Client
}

func (c *Client) Get(o *Options) (*WeatherData, error) {
	req, err := http.NewRequest("GET", c.url(o), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var wd WeatherData
	if err = json.NewDecoder(res.Body).Decode(&wd); err != nil {
		return nil, err
	}

	return &wd, nil
}

type Options struct {
	Latitude          float64
	Longitude         float64
	TemperatureUnit   string        // Default "celsius"
	WindspeedUnit     string        // Default "kmh",
	PrecipitationUnit string        // Default "mm"
	Timezone          time.Location // Default "UTC"
	PastDays          int           // Default 0
	ForcastDays       int
	Start             time.Time
	End               time.Time
	HourlyMetrics     *HourlyMetrics
	DailyMetrics      *DailyMetrics
	CurrentMetrics    *CurrentMetrics
}

type HourlyMetrics struct {
	Temperature2m            bool `json:"temperature_2m"`
	RelativeHumidity2m       bool `json:"relative_humidity_2m"`
	DewPoint2m               bool `json:"dew_point_2m"`
	ApparentTemperature      bool `json:"apparent_temperature"`
	PrecipitationProbability bool `json:"precipitation_probability"`
	Precipitation            bool `json:"precipitation"`
	Rain                     bool `json:"rain"`
	Showers                  bool `json:"showers"`
	Snowfall                 bool `json:"snowfall"`
	SnowDepth                bool `json:"snow_depth"`
	WeatherCode              bool `json:"weather_code"`
	PressureMsl              bool `json:"pressure_msl"`
	SurfacePressure          bool `json:"surface_pressure"`
	CloudCover               bool `json:"cloud_cover"`
	CloudCoverLow            bool `json:"cloud_cover_low"`
	CloudCoverMid            bool `json:"cloud_cover_mid"`
	CloudCoverHigh           bool `json:"cloud_cover_high"`
	Evapotranspiration       bool `json:"evapotranspiration"`
	Visibility               bool `json:"visibility"`
	Et0FaoEvapotranspiration bool `json:"et0_fao_evapotranspiration"`
	VapourPressureDeficit    bool `json:"vapour_pressure_deficit"`
	WindSpeed10m             bool `json:"wind_speed_10m"`
	WindSpeed80m             bool `json:"wind_speed_80m"`
	WindSpeed120m            bool `json:"wind_speed_120m"`
	WindSpeed180m            bool `json:"wind_speed_180m"`
	WindDirection10m         bool `json:"wind_direction_10m"`
	WindDirection80m         bool `json:"wind_direction_80m"`
	WindDirection120m        bool `json:"wind_direction_120m"`
	WindDirection180m        bool `json:"wind_direction_180m"`
	WindGusts10m             bool `json:"wind_gusts_10m"`
	Temperature80m           bool `json:"temperature_80m"`
	Temperature120m          bool `json:"temperature_120m"`
	Temperature180m          bool `json:"temperature_180m"`
	SoilTemperature0cm       bool `json:"soil_temperature_0cm"`
	SoilTemperature6cm       bool `json:"soil_temperature_6cm"`
	SoilTemperature18cm      bool `json:"soil_temperature_18cm"`
	SoilTemperature54cm      bool `json:"soil_temperature_54cm"`
	SoilMoisture0To1cm       bool `json:"soil_moisture_0_to_1cm"`
	SoilMoisture1To3cm       bool `json:"soil_moisture_1_to_3cm"`
	SoilMoisture9To27cm      bool `json:"soil_moisture_9_to_27cm"`
	SoilMoisture3To9cm       bool `json:"soil_moisture_3_to_9cm"`
}

// encodeMetrics uses reflection to create a comma-separated string of json tags for a struct's boolean fields that are true.
func encodeMetrics(s interface{}) string {
	var metrics []string
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Bool && field.Bool() {
			metrics = append(metrics, t.Field(i).Tag.Get("json"))
		}
	}
	return strings.Join(metrics, ",")
}

func (h *HourlyMetrics) encode() string {
	return encodeMetrics(h)
}

// CurrentMetrics specifies which current weather variables to retrieve.
type CurrentMetrics struct {
	Temperature2m       bool `json:"temperature_2m"`
	RelativeHumidity2m  bool `json:"relative_humidity_2m"`
	IsDay               bool `json:"is_day"`
	ApparentTemperature bool `json:"apparent_temperature"`
	Precipitation       bool `json:"precipitation"`
	Rain                bool `json:"rain"`
	Showers             bool `json:"showers"`
	Snowfall            bool `json:"snowfall"`
	WeatherCode         bool `json:"weather_code"`
	CloudCover          bool `json:"cloud_cover"`
	PressureMsl         bool `json:"pressure_msl"`
	SurfacePressure     bool `json:"surface_pressure"`
	WindSpeed10m        bool `json:"wind_speed_10m"`
	WindDirection10m    bool `json:"wind_direction_10m"`
	WindGusts10m        bool `json:"wind_gusts_10m"`
}

func (c *CurrentMetrics) encode() string {
	return encodeMetrics(c)
}

// DailyMetrics specifies which daily weather variables to retrieve.
type DailyMetrics struct {
	WeatherCode                 bool `json:"weather_code"`
	Temperature2mMax            bool `json:"temperature_2m_max"`
	Temperature2mMin            bool `json:"temperature_2m_min"`
	ApparentTemperatureMax      bool `json:"apparent_temperature_max"`
	ApparentTemperatureMin      bool `json:"apparent_temperature_min"`
	Sunrise                     bool `json:"sunrise"`
	Sunset                      bool `json:"sunset"`
	SunshineDuration            bool `json:"sunshine_duration"`
	DaylightDuration            bool `json:"daylight_duration"`
	UvIndexMax                  bool `json:"uv_index_max"`
	UvIndexClearSkyMax          bool `json:"uv_index_clear_sky_max"`
	RainSum                     bool `json:"rain_sum"`
	ShowersSum                  bool `json:"showers_sum"`
	SnowfallSum                 bool `json:"snowfall_sum"`
	PrecipitationSum            bool `json:"precipitation_sum"`
	PrecipitationHours          bool `json:"precipitation_hours"`
	PrecipitationProbabilityMax bool `json:"precipitation_probability_max"`
	WindSpeed10mMax             bool `json:"wind_speed_10m_max"`
	WindGusts10mMax             bool `json:"wind_gusts_10m_max"`
	WindDirection10mDominant    bool `json:"wind_direction_10m_dominant"`
	ShortwaveRadiationSum       bool `json:"shortwave_radiation_sum"`
	Et0FaoEvapotranspiration    bool `json:"et0_fao_evapotranspiration"`
}

func (d *DailyMetrics) encode() string {
	return encodeMetrics(d)
}

func New() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		UserAgent:  DefaultUserAgent,
	}
}

func NewWithKey(key string) *Client {
	return &Client{
		apiKey:     key,
		HTTPClient: http.DefaultClient,
		UserAgent:  DefaultUserAgent,
	}
}

func (c *Client) url(o *Options) string {
	path := "/v1/forecast"
	prefixes := []string{}
	host := baseHost

	if (!o.Start.IsZero()) && time.Since(o.Start) > 7*24*time.Hour {
		prefixes = append(prefixes, "archive-")
		path = "/v1/archive"
	}

	if c.apiKey != "" {
		prefixes = append(prefixes, "customer-")
	}

	for _, prefix := range prefixes {
		host = prefix + host
	}

	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}

	q := newOrderedQuery()

	if c.apiKey != "" {
		q.Set("apikey", c.apiKey)
	}

	q.Set("latitude", fmt.Sprintf("%v", o.Latitude))
	q.Set("longitude", fmt.Sprintf("%v", o.Longitude))

	if o.TemperatureUnit != "" {
		q.Set("temperature_unit", o.TemperatureUnit)
	}

	if o.WindspeedUnit != "" {
		q.Set("windspeed_unit", o.WindspeedUnit)
	}

	if o.PrecipitationUnit != "" {
		q.Set("precipitation_unit", o.PrecipitationUnit)
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
	} else if o.ForcastDays > 0 {
		// Default to weathercode if forecast days are requested but no specific daily metrics
		q.Set("daily", "weathercode")
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

type orderedQuery struct {
	keys   []string
	values map[string]string
}

func (o *orderedQuery) Set(key, value string) {
	o.keys = append(o.keys, key)
	o.values[key] = value
}

func (o *orderedQuery) Encode() string {
	tmp := []string{}

	for _, key := range o.keys {
		tmp = append(tmp, fmt.Sprintf("%s=%s", key, o.values[key]))
	}

	return strings.Join(tmp, "&")
}

func newOrderedQuery() *orderedQuery {
	return &orderedQuery{
		values: map[string]string{},
		keys:   []string{},
	}
}

// WeatherData is the main struct that holds all the forecast data.
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
