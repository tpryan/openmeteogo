package openmateogo

import "net/http"

const baseHost = "api.open-meteo.com/v1/forecast"
const histHost = "archive-api.open-meteo.com/v1/archive"
const apiPrefix = "server-"

const DefaultUserAgent = "OpenMeteoGo-Client"

type Client struct {
	apiKey     string
	UserAgent  string
	HTTPClient *http.Client
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
