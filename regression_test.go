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

// Regression tests for defects that the rest of the suite executed without
// noticing. Each one fails against the behaviour that preceded it.

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// capturingTransport answers every request from canned data and keeps the
// request it was given. Replacing the transport keeps these tests off the
// network and away from DNS, which matters for the commercial host names.
type capturingTransport struct {
	gotReq *http.Request
	status int
	body   string
	header http.Header
}

func (c *capturingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	c.gotReq = r.Clone(r.Context())

	h := c.header
	if h == nil {
		h = make(http.Header)
	}

	status := c.status
	if status == 0 {
		status = http.StatusOK
	}

	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(c.body)),
		Header:     h,
		Request:    r,
	}, nil
}

// clientWithResponse returns a client whose transport always answers with the
// given status, body and headers.
func clientWithResponse(status int, body string, header http.Header) (*Client, *capturingTransport) {
	tr := &capturingTransport{status: status, body: body, header: header}
	c := NewClient()
	c.HTTPClient = &http.Client{Transport: tr}

	return c, tr
}

// --- Timezone --------------------------------------------------------------

func TestTimezone_ReachesTheURL(t *testing.T) {
	c := NewClient()

	la, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)

	t.Run("loaded location", func(t *testing.T) {
		assert.Contains(t, c.url(NewOptionsBuilder().Timezone(la).Build()),
			"timezone=America%2FLos_Angeles")
	})

	t.Run("UTC", func(t *testing.T) {
		assert.Contains(t, c.url(NewOptionsBuilder().Timezone(time.UTC).Build()),
			"timezone=UTC")
	})

	t.Run("nil leaves the parameter off", func(t *testing.T) {
		assert.NotContains(t, c.url(NewOptionsBuilder().Build()), "timezone=")
	})

	t.Run("the builder keeps the caller's location rather than a copy", func(t *testing.T) {
		// The copy was the defect. time.Local is a sentinel the time package
		// recognises by address, so a copy of it reports an empty name until
		// the package happens to have initialised, and the timezone silently
		// fell out of the request depending on what ran first.
		opts := NewOptionsBuilder().Timezone(time.Local).Build()
		assert.Same(t, time.Local, opts.Timezone)
	})
}

// --- Endpoint routing ------------------------------------------------------

func TestURL_ArchiveOnlyWhenTheWholeRangeIsHistorical(t *testing.T) {
	c := NewClient()
	now := time.Now()

	t.Run("a fully historical range goes to the archive", func(t *testing.T) {
		o := NewOptionsBuilder().
			Start(now.AddDate(0, 0, -60)).
			End(now.AddDate(0, 0, -30)).
			Build()
		assert.Contains(t, c.url(o), "archive-api.open-meteo.com/v1/archive")
	})

	t.Run("an old start with no end goes to the archive", func(t *testing.T) {
		o := NewOptionsBuilder().Start(now.AddDate(0, 0, -60)).Build()
		assert.Contains(t, c.url(o), "archive-api.open-meteo.com/v1/archive")
	})

	t.Run("a range reaching into the forecast window stays on forecast", func(t *testing.T) {
		// Routing this to the archive silently drops the forecast half of
		// the range and returns a short series instead of an error.
		o := NewOptionsBuilder().
			Start(now.AddDate(0, 0, -30)).
			End(now.AddDate(0, 0, 3)).
			Build()

		u := c.url(o)
		assert.NotContains(t, u, "archive")
		assert.Contains(t, u, "https://api.open-meteo.com/v1/forecast")
	})
}

// --- API key handling ------------------------------------------------------

func TestGet_APIKeyTravelsInTheHeaderNotTheQuery(t *testing.T) {
	tr := &capturingTransport{body: `{"latitude":52.52}`}
	c := NewClientWithKey("secret-key-value")
	c.HTTPClient = &http.Client{Transport: tr}

	_, err := c.Get(context.Background(), NewOptionsBuilder().Latitude(52.52).Build())
	require.NoError(t, err)
	require.NotNil(t, tr.gotReq)

	assert.Equal(t, "secret-key-value", tr.gotReq.Header.Get("X-Api-Key"))

	// A transport failure returns a *url.Error holding the URL verbatim, so
	// a key in the query string ends up in the caller's logs.
	assert.Empty(t, tr.gotReq.URL.Query().Get("apikey"))
	assert.NotContains(t, tr.gotReq.URL.String(), "secret-key-value")

	// The key still selects the commercial host.
	assert.Equal(t, "customer-api.open-meteo.com", tr.gotReq.URL.Host)
}

// --- Error reporting -------------------------------------------------------

func TestAPIError_CarriesTheReasonFromTheBody(t *testing.T) {
	bodies := map[string]string{
		"boolean error flag": `{"error":true,"reason":"Invalid timezone"}`,
		"quoted error flag":  `{"error":"true","reason":"Invalid timezone"}`,
	}

	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			c, _ := clientWithResponse(http.StatusBadRequest, body, nil)

			_, err := c.Get(context.Background(), NewOptionsBuilder().Build())
			require.Error(t, err)

			var apiErr *APIError
			require.ErrorAs(t, err, &apiErr)
			assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
			assert.Equal(t, "Invalid timezone", apiErr.Reason)
			assert.Contains(t, apiErr.Error(), "Invalid timezone")
		})
	}
}

func TestAPIError_RetryAfterIsParsed(t *testing.T) {
	t.Run("delay in seconds", func(t *testing.T) {
		h := http.Header{}
		h.Set("Retry-After", "30")
		c, _ := clientWithResponse(http.StatusTooManyRequests, `{"reason":"Rate limited"}`, h)

		_, err := c.Get(context.Background(), NewOptionsBuilder().Build())

		var apiErr *APIError
		require.ErrorAs(t, err, &apiErr)
		assert.True(t, apiErr.RateLimited())
		assert.Equal(t, 30*time.Second, apiErr.RetryAfter)
		assert.Contains(t, apiErr.Error(), "retry after 30s")
	})

	t.Run("absent header leaves it zero", func(t *testing.T) {
		c, _ := clientWithResponse(http.StatusTooManyRequests, `{"reason":"Rate limited"}`, nil)

		_, err := c.Get(context.Background(), NewOptionsBuilder().Build())

		var apiErr *APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Zero(t, apiErr.RetryAfter)
	})
}

func TestAPIError_NonJSONBodyIsKeptAsTheReason(t *testing.T) {
	c, _ := clientWithResponse(http.StatusBadGateway, "<html>gateway fell over</html>", nil)

	_, err := c.Get(context.Background(), NewOptionsBuilder().Build())

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusBadGateway, apiErr.StatusCode)
	assert.Contains(t, apiErr.Reason, "gateway fell over")
}

// --- Context ---------------------------------------------------------------

func TestGet_ContextCancellationAbortsTheRequest(t *testing.T) {
	release := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
	}))
	defer server.Close()
	defer close(release)

	c := NewClient()
	c.HTTPClient = server.Client()

	parts := strings.Split(server.URL, "://")
	c.scheme, c.host = parts[0], parts[1]

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := c.Get(ctx, NewOptionsBuilder().Build())
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestGet_NilOptionsIsAnErrorNotAPanic(t *testing.T) {
	c := NewClient()

	assert.NotPanics(t, func() {
		_, err := c.Get(context.Background(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nil Options")
	})
}

// --- Client construction ---------------------------------------------------

func TestNewClient_HasATimeoutAndIsNotTheGlobalDefault(t *testing.T) {
	c := NewClient()
	require.NotNil(t, c.HTTPClient)

	// Sharing http.DefaultClient means a caller setting a timeout on it
	// reaches into every other user of the default client in the process.
	assert.NotSame(t, http.DefaultClient, c.HTTPClient)
	assert.Equal(t, defaultTimeout, c.HTTPClient.Timeout)
}

func TestZeroValueClientFallsBackToDefaults(t *testing.T) {
	var c Client

	// Client exports HTTPClient and UserAgent but not scheme or host, so a
	// struct literal is a natural thing to write and used to produce a
	// malformed URL and then nil-panic in Do.
	u := c.url(NewOptionsBuilder().Latitude(1).Longitude(2).Build())
	assert.Equal(t, "https://api.open-meteo.com/v1/forecast?latitude=1&longitude=2", u)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // fail before dialling, so the test stays off the network

	assert.NotPanics(t, func() {
		_, err := c.Get(ctx, NewOptionsBuilder().Build())
		assert.Error(t, err)
	})
}

// --- Metrics ---------------------------------------------------------------

func TestNewMetrics_UnknownTypeIsAnError(t *testing.T) {
	// This used to leave the allowed list nil, so a call passing no metrics
	// returned an empty Metrics and a nil error, accepting the typo.
	got, err := NewMetrics("Hourly")
	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "unknown metric type")

	got, err = NewMetrics("hourly-ish", Temperature2m)
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestNewMetrics_SeasonalMetricsAreValidated(t *testing.T) {
	got, err := NewMetrics("weekly", Temperature2mMean, PrecipitationAnomaly)
	require.NoError(t, err)
	assert.Equal(t, Metrics{Temperature2mMean, PrecipitationAnomaly}, got)

	got, err = NewMetrics("monthly", SoilMoisture0To10cmMean)
	require.NoError(t, err)
	assert.Equal(t, Metrics{SoilMoisture0To10cmMean}, got)

	// A forecast metric is not a seasonal one. Both arms used to return the
	// caller's input unchecked.
	_, err = NewMetrics("weekly", Temperature2m)
	require.Error(t, err)
}

func TestMetrics_EncodeDedupesAndSkipsEmpty(t *testing.T) {
	m := Metrics{Temperature2m, Temperature2m, "", WindSpeed10m}
	assert.Equal(t, "temperature_2m,wind_speed_10m", m.encode())

	assert.Empty(t, Metrics(nil).encode())
	assert.Empty(t, Metrics{}.encode())
}

// --- Coordinate formatting -------------------------------------------------

func TestFormatCoord_NeverUsesScientificNotation(t *testing.T) {
	// The %v verb renders small magnitudes as 1e-07, which the API rejects.
	assert.Equal(t, "0.0000001", formatCoord(0.0000001))
	assert.Equal(t, "52.52", formatCoord(52.52))
	assert.Equal(t, "-122.4194", formatCoord(-122.4194))
	assert.Equal(t, "0", formatCoord(0))
}
