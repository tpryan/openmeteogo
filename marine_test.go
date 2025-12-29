package openmeteogo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURL_Marine(t *testing.T) {
	tests := map[string]struct {
		client  *Client
		options Options
		want    string
	}{
		"basic marine via flag": {
			client:  NewClient(),
			options: *NewOptionsBuilder().Marine(true).Latitude(52.52).Longitude(13.41).Build(),
			want:    "https://marine-api.open-meteo.com/v1/marine?latitude=52.52&longitude=13.41",
		},
		"marine with hourly metrics": {
			client:  NewClient(),
			options: *NewOptionsBuilder().Marine(true).HourlyMetrics(Metrics{WaveHeight}).Latitude(0).Longitude(0).Build(),
			want:    "https://marine-api.open-meteo.com/v1/marine?hourly=wave_height&latitude=0&longitude=0",
		},
		"marine with daily metrics": {
			client:  NewClient(),
			options: *NewOptionsBuilder().Marine(true).DailyMetrics(Metrics{WaveHeightMax}).Latitude(0).Longitude(0).Build(),
			want:    "https://marine-api.open-meteo.com/v1/marine?daily=wave_height_max&latitude=0&longitude=0",
		},
		"with api key": {
			client:  NewClientWithKey("testkey"),
			options: *NewOptionsBuilder().Marine(true).Latitude(0).Longitude(0).Build(),
			want:    "https://customer-marine-api.open-meteo.com/v1/marine?apikey=testkey&latitude=0&longitude=0",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.client.url(&tc.options)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestClient_Get_Marine(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/v1/marine" {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{
			"latitude": 52.52,
			"longitude": 13.41,
			"hourly_units": {
				"time": "iso8601",
				"wave_height": "m"
			},
			"hourly": {
				"time": ["2025-01-01T00:00"],
				"wave_height": [2.5]
			}
		}`))
	}))
	defer server.Close()

	client := NewClient()
	client.HTTPClient = server.Client()
	urlParts := strings.Split(server.URL, "://")
	client.scheme = urlParts[0]
	client.host = urlParts[1]
	client.marineHost = urlParts[1]

	opts := NewOptionsBuilder().
		Marine(true).
		Latitude(52.52).
		Longitude(13.41).
		HourlyMetrics(Metrics{WaveHeight}).
		Build()

	wd, err := client.Get(opts)
	require.NoError(t, err)
	assert.NotNil(t, wd)
	assert.Equal(t, 52.52, wd.Latitude)
	assert.Equal(t, "m", wd.HourlyUnits.WaveHeight)
	assert.Equal(t, 1, len(wd.Hourly.Time))
	assert.Equal(t, 2.5, wd.Hourly.WaveHeight[0])
}
