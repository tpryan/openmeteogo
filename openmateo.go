package openmateogo

import (
	"fmt"
	"net/http"
	"net/url"
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

type Options struct {
	Latitude          float64
	Longitude         float64
	TemperatureUnit   string         // Default "celsius"
	WindspeedUnit     string         // Default "kmh",
	PrecipitationUnit string         // Default "mm"
	Timezone          *time.Location // Default "UTC"
	PastDays          int            // Default 0
	ForcastDays       int
	Start             time.Time
	End               time.Time
	// HourlyMetrics     []string // Lists required hourly metrics, see https://open-meteo.com/en/docs for valid metrics
	// DailyMetrics      []string // Lists required daily metrics, see https://open-meteo.com/en/docs for valid metrics
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

	q := u.Query()
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

	if o.Timezone == nil {
		q.Set("timezone", o.Timezone.String())
	}

	if o.PastDays > 0 {
		q.Set("past_days", fmt.Sprintf("%v", o.PastDays))
	}

	if o.ForcastDays > 0 {
		q.Set("daily", "weathercode")
		q.Set("daily_units", "weathercode")
		q.Set("forecast_days", fmt.Sprintf("%v", o.ForcastDays))
	}

	u.RawQuery = q.Encode()

	return u.String()
}
