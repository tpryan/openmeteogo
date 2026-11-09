package openmeteogo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURL_Seasonal(t *testing.T) {
	tests := map[string]struct {
		client  *Client
		options Options
		want    string
	}{
		"basic seasonal via models": {
			client:  NewClient(),
			options: *NewOptionsBuilder().Models([]string{"ecmwf_seas5"}).Latitude(52.52).Longitude(13.41).Build(),
			want:    "https://seasonal-api.open-meteo.com/v1/seasonal?latitude=52.52&longitude=13.41&models=ecmwf_seas5",
		},
		"basic seasonal via flag": {
			client:  NewClient(),
			options: *NewOptionsBuilder().Seasonal(true).Latitude(52.52).Longitude(13.41).Build(),
			want:    "https://seasonal-api.open-meteo.com/v1/seasonal?latitude=52.52&longitude=13.41",
		},
		"with models": {
			client:  NewClient(),
			options: *NewOptionsBuilder().Models([]string{"ecmwf_seas5", "cfs"}).Latitude(0).Longitude(0).Build(),
			want:    "https://seasonal-api.open-meteo.com/v1/seasonal?latitude=0&longitude=0&models=ecmwf_seas5%2Ccfs",
		},
		"weekly metrics": {
			client:  NewClient(),
			options: *NewOptionsBuilder().WeeklyMetrics(Metrics{Temperature2mMean}).Latitude(0).Longitude(0).Build(),
			want:    "https://seasonal-api.open-meteo.com/v1/seasonal?latitude=0&longitude=0&weekly=temperature_2m_mean",
		},
		"monthly metrics": {
			client:  NewClient(),
			options: *NewOptionsBuilder().MonthlyMetrics(Metrics{PrecipitationMean}).Latitude(0).Longitude(0).Build(),
			want:    "https://seasonal-api.open-meteo.com/v1/seasonal?latitude=0&longitude=0&monthly=precipitation_mean",
		},
		"complex seasonal": {
			client: NewClient(),
			options: *NewOptionsBuilder().
				Latitude(10).
				Longitude(20).
				Models([]string{"modelA"}).
				WeeklyMetrics(Metrics{Temperature2mMean}).
				MonthlyMetrics(Metrics{PrecipitationMean}).
				Build(),
			want: "https://seasonal-api.open-meteo.com/v1/seasonal?latitude=10&longitude=20&models=modelA&monthly=precipitation_mean&weekly=temperature_2m_mean",
		},
		"with api key": {
			client:  NewClientWithKey("testkey"),
			options: *NewOptionsBuilder().Models([]string{"modelA"}).Latitude(0).Longitude(0).Build(),
			want:    "https://customer-seasonal-api.open-meteo.com/v1/seasonal?apikey=testkey&latitude=0&longitude=0&models=modelA",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.client.url(&tc.options)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestClient_Get_Seasonal(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Verify request path
		if req.URL.Path != "/v1/seasonal" {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		// Send minimal valid response
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{
			"latitude": 52.52,
			"longitude": 13.41,
			"weekly_units": {
				"time": "iso8601",
				"temperature_2m_mean": "°C"
			},
			"weekly": {
				"time": ["2025-01-01"],
				"temperature_2m_mean": [10.5]
			}
		}`))
	}))
	defer server.Close()

	client := NewClient()
	client.HTTPClient = server.Client()
	
	// Override hosts to point to mock server
	urlParts := strings.Split(server.URL, "://")
	client.scheme = urlParts[0]
	client.host = urlParts[1]
	client.seasonalHost = urlParts[1]

	opts := NewOptionsBuilder().
		Latitude(52.52).
		Longitude(13.41).
		WeeklyMetrics(Metrics{Temperature2mMean}).
		Build()

	wd, err := client.Get(opts)
	require.NoError(t, err)
	assert.NotNil(t, wd)
	assert.Equal(t, 52.52, wd.Latitude)
	assert.Equal(t, 13.41, wd.Longitude)
	assert.NotNil(t, wd.Weekly)
	assert.Equal(t, "°C", wd.WeeklyUnits.Temperature2mMean)
	assert.Equal(t, 1, len(wd.Weekly.Time))
	assert.Equal(t, 10.5, wd.Weekly.Temperature2mMean[0])
}

func TestClient_Get_Seasonal_Integration(t *testing.T) {
	// This test uses the mock server but simulates a more complete response
	// to ensure all fields are mapped correctly if JSON is full.
	
	responseJSON := `{
		"latitude": 10.0,
		"longitude": 20.0,
		"generationtime_ms": 1.23,
		"utc_offset_seconds": 0,
		"timezone": "GMT",
		"timezone_abbreviation": "GMT",
		"elevation": 100.0,
		"weekly_units": {
			"time": "iso8601",
			"temperature_2m_mean": "°C",
			"temperature_2m_anomaly": "°C",
			"precipitation_mean": "mm",
			"precipitation_anomaly": "mm"
		},
		"weekly": {
			"time": ["2025-01-01", "2025-01-08"],
			"temperature_2m_mean": [10.0, 12.0],
			"temperature_2m_anomaly": [1.0, 0.5],
			"precipitation_mean": [5.0, 2.0],
			"precipitation_anomaly": [-1.0, 0.0]
		},
		"monthly_units": {
			"time": "iso8601",
			"temperature_2m_mean": "°C"
		},
		"monthly": {
			"time": ["2025-01-01"],
			"temperature_2m_mean": [11.0]
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(responseJSON))
	}))
	defer server.Close()

	client := NewClient()
	client.HTTPClient = server.Client()
	urlParts := strings.Split(server.URL, "://")
	client.scheme = urlParts[0]
	client.seasonalHost = urlParts[1]

	opts := NewOptionsBuilder().
		WeeklyMetrics(Metrics{Temperature2mMean}). // Triggers seasonal logic
		Build()
	wd, err := client.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, 10.0, wd.Latitude)
	assert.Equal(t, 2, len(wd.Weekly.Time))
	assert.Equal(t, 12.0, wd.Weekly.Temperature2mMean[1])
	assert.Equal(t, 1.0, wd.Weekly.Temperature2mAnomaly[0])
	assert.Equal(t, 1, len(wd.Monthly.Time))
	assert.Equal(t, 11.0, wd.Monthly.Temperature2mMean[0])
}
