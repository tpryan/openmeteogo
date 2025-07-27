package openmateo

import (
	"reflect"
	"strings"
	"time"
)

type Options struct {
	Latitude          float64
	Longitude         float64
	TemperatureUnit   TemperatureUnit   // Default "celsius"
	WindspeedUnit     WindSpeedUnit     // Default "kmh",
	PrecipitationUnit PrecipitationUnit // Default "mm"
	Timezone          time.Location     // Default "UTC"
	PastDays          int               // Default 0
	ForcastDays       int
	Start             time.Time
	End               time.Time
	HourlyMetrics     *HourlyMetrics
	DailyMetrics      *DailyMetrics
	CurrentMetrics    *CurrentMetrics
}

type OptionsBuilder struct {
	options *Options
}

func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{options: &Options{}}
}

func (b *OptionsBuilder) Latitude(lat float64) *OptionsBuilder {
	b.options.Latitude = lat
	return b
}

func (b *OptionsBuilder) Longitude(lon float64) *OptionsBuilder {
	b.options.Longitude = lon
	return b
}

func (b *OptionsBuilder) TemperatureUnit(unit TemperatureUnit) *OptionsBuilder {
	b.options.TemperatureUnit = unit
	return b
}

func (b *OptionsBuilder) WindspeedUnit(unit WindSpeedUnit) *OptionsBuilder {
	b.options.WindspeedUnit = unit
	return b
}

func (b *OptionsBuilder) PrecipitationUnit(unit PrecipitationUnit) *OptionsBuilder {
	b.options.PrecipitationUnit = unit
	return b
}

func (b *OptionsBuilder) Timezone(tz time.Location) *OptionsBuilder {
	b.options.Timezone = tz
	return b
}

func (b *OptionsBuilder) PastDays(days int) *OptionsBuilder {
	b.options.PastDays = days
	return b
}

func (b *OptionsBuilder) ForcastDays(days int) *OptionsBuilder {
	b.options.ForcastDays = days
	return b
}

func (b *OptionsBuilder) Start(start time.Time) *OptionsBuilder {
	b.options.Start = start
	return b
}

func (b *OptionsBuilder) End(end time.Time) *OptionsBuilder {
	b.options.End = end
	return b
}

func (b *OptionsBuilder) HourlyMetrics(metrics *HourlyMetrics) *OptionsBuilder {
	b.options.HourlyMetrics = metrics
	return b
}

func (b *OptionsBuilder) DailyMetrics(metrics *DailyMetrics) *OptionsBuilder {
	b.options.DailyMetrics = metrics
	return b
}

func (b *OptionsBuilder) CurrentMetrics(metrics *CurrentMetrics) *OptionsBuilder {
	b.options.CurrentMetrics = metrics
	return b
}

func (b *OptionsBuilder) Build() *Options {
	return b.options
}

// HourlyMetrics specifies which hourly weather variables to retrieve.

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
