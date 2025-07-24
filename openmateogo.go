package openmateogo

import (
	"net/http"
	"time"
)

const baseHost = "api.open-meteo.com/v1/forecast"
const histHost = "archive-api.open-meteo.com/v1/archive"
const apiPrefix = "server-"

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
	TemperatureUnit   string        // Default "celsius"
	WindspeedUnit     string        // Default "kmh",
	PrecipitationUnit string        // Default "mm"
	Timezone          time.Location // Default "UTC"
	PastDays          int           // Default 0
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
	return ""
}
