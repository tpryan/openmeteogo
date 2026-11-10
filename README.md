# **OpenMeteo Go Client**

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/tpryan/openmeteogo?status.svg)](https://godoc.org/github.com/tpryan/openmeteogo)

A simple and easy-to-use Go client for the 
[Open-Meteo API](https://open-meteo.com/). This library provides a convenient
way to fetch current, historical, and forecasted weather data.

I built this mainly because I couldn't find historical support in the other 
go libraries for Open Meteo, and because it made for a fun weekend project.


## **Features**

* Fetch **current**, **hourly**, and **daily** weather data.  
* Access **historical weather data** for specific date ranges.  
* Fetch **marine data** (waves, swell, ocean currents), with JSON `null` kept distinct from a real `0`.  
* Fluent **Options Builder** for easily constructing complex queries.  
* Support for multiple **units** (Celsius/Fahrenheit, km/h/mph, etc.).  
* Select exactly which weather metrics you need.  
* Helper function to get human-readable descriptions from WMO weather codes.  
* Supports both the free and commercial (API key) Open-Meteo endpoints.

## **Installation**

To install the library, use go get:

```bash
go get github.com/tpryan/openmeteogo/v3
```

> **Upgrading from v2?** `Get` now takes a `context.Context`, `Timezone` takes a
> `*time.Location`, `ForcastDays` is spelled `ForecastDays`, and non-200
> responses return an `*APIError`. See [Handling Errors](#handling-errors).

## **Usage**

Here is a quick example of how to get the current weather and today's forecast
for San Francisco.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tpryan/openmeteogo/v3"
)

func main() {
	// 1. Create a new client
	c := openmeteogo.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 2. Set the options for the weather data you want.
	// Use the fluent builder to easily configure your request.
	opts := openmeteogo.NewOptionsBuilder().
		Latitude(37.7749).
		Longitude(-122.4194).
		TemperatureUnit(openmeteogo.Fahrenheit).
		ForecastDays(1).
		CurrentMetrics(openmeteogo.Metrics{
			openmeteogo.Temperature2m,
			openmeteogo.WeatherCode,
		}).
		DailyMetrics(openmeteogo.Metrics{
			openmeteogo.Temperature2mMax,
			openmeteogo.Temperature2mMin,
			openmeteogo.WeatherCode,
		}).
		Build()

	// 3. Make the API call
	w, err := c.Get(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to get weather data: %v", err)
	}

	// 4. Use the data!
	fmt.Printf("Current Weather in San Francisco (Lat: %.2f, Lon: %.2f)\n", w.Latitude, w.Longitude)
	fmt.Printf("Time: %s\n", w.Current.Time)
	fmt.Printf("Temperature: %.2f%s\n", w.Current.Temperature2m, w.CurrentUnits.Temperature2m)
	fmt.Printf("Weather: %s\n\n", openmeteogo.DescribeCode(w.Current.WeatherCode))

	// Print the daily forecast for today
	if len(w.Daily.Time) > 0 {
		fmt.Printf("Forecast for %s:\n", w.Daily.Time[0])
		fmt.Printf("  Max Temperature: %.2f%s\n", w.Daily.Temperature2mMax[0], w.DailyUnits.Temperature2mMax)
		fmt.Printf("  Min Temperature: %.2f%s\n", w.Daily.Temperature2mMin[0], w.DailyUnits.Temperature2mMin)
		fmt.Printf("  Weather: %s\n", openmeteogo.DescribeCode(w.Daily.WeatherCode[0]))
	}
}
```

### **Fetching Historical Data**

To get historical data, simply provide a Start and End date. The client will automatically use the correct API endpoint.


```go

    // Fetch historical data for a specific week
    pastOpts := openmeteogo.NewOptionsBuilder().
		Latitude(37.7749).
        Longitude(-122.4194).
		TemperatureUnit(openmeteogo.Fahrenheit).
		DailyMetrics(openmeteogo.Metrics{
			openmeteogo.Temperature2mMax,
			openmeteogo.Temperature2mMin,
			openmeteogo.WeatherCode,
		}).
		Start(time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)).
        End(time.Date(2023, 7, 7, 0, 0, 0, 0, time.UTC)).
		Build()

	wp, err := c.Get(ctx, pastOpts)
	if err != nil {
		log.Fatalf("Failed to get historical weather data: %v", err)
	}

    // Process the historical daily data...
    for i, date := range wp.Daily.Time {
        fmt.Printf("Weather for %s:\n", date)
		fmt.Printf("  Max Temp: %.2f%s\n", wp.Daily.Temperature2mMax[i], wp.DailyUnits.Temperature2mMax)
		fmt.Printf("  Min Temp: %.2f%s\n", wp.Daily.Temperature2mMin[i], wp.DailyUnits.Temperature2mMin)
		fmt.Printf("  Weather: %s\n", openmeteogo.DescribeCode(wp.Daily.WeatherCode[i]))
    }
```    

### **Fetching Marine Data**

Set `Marine(true)` to route the request to the marine endpoint, which serves
wave, swell and ocean-current data.

The marine model only covers open water. A position outside its domain — a
marina well up an estuary, say — does **not** return an error: it returns HTTP
200 with a full-length series of JSON `null`. Marine values are therefore
`*float64` rather than `float64`, so `nil` means *no value here* and `0` means a
real reading of zero. Check coverage with `HasAny` before reading a series.

```go
	opts := openmeteogo.NewOptionsBuilder().
		Latitude(38.0447).
		Longitude(-122.1577).
		Marine(true).
		DailyMetrics(openmeteogo.Metrics{
			openmeteogo.WaveHeightMax,
			openmeteogo.WavePeriodMax,
		}).
		Build()

	wd, err := c.Get(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to get marine data: %v", err)
	}

	if !openmeteogo.HasAny(wd.Daily.WaveHeightMax) {
		log.Fatal("no wave data: this position is outside the marine model's domain")
	}

	for i, date := range wd.Daily.Time {
		h := wd.Daily.WaveHeightMax[i]
		if h == nil {
			fmt.Printf("%s: no wave height\n", date)
			continue
		}
		fmt.Printf("%s: max wave height %.2f%s\n", date, *h, wd.DailyUnits.WaveHeightMax)
	}
```

Three helpers keep callers from reinventing the nil handling:

| Helper | Description |
| :---- | :---- |
| `Value(vals []*float64) (float64, bool)` | First non-nil entry, and whether one existed. |
| `HasAny(vals []*float64) bool` | Whether any entry is non-nil. False across a marine series means the position is out of domain. |
| `Floats(vals []*float64) []float64` | Flattens to `[]float64`, substituting zero for nil. Only safe once `HasAny` is true. |

Validate marine metrics with `NewMetrics` using the `marine_hourly`,
`marine_daily`, or `marine_current` metric type.

### **Handling Errors**

A rejected request comes back as an `*APIError`. Open-Meteo explains failures
in the response body, and `Reason` carries that text, so you get "Invalid
timezone" rather than only a status code.

```go
	w, err := c.Get(ctx, opts)
	if err != nil {
		var apiErr *openmeteogo.APIError
		if errors.As(err, &apiErr) {
			if apiErr.RateLimited() {
				log.Fatalf("rate limited, retry after %s", apiErr.RetryAfter)
			}
			log.Fatalf("http %d: %s", apiErr.StatusCode, apiErr.Reason)
		}
		log.Fatalf("request failed: %v", err)
	}
```

`RetryAfter` is populated from the `Retry-After` header when the API sends one,
and is zero otherwise.

## **Options**

The OptionsBuilder provides a simple way to configure your request.

| Method | Description | Example |
| :---- | :---- | :---- |
| Latitude() | Set the geographical latitude. | .Latitude(37.7749) |
| Longitude() | Set the geographical longitude. | .Longitude(-122.4194) |
| TemperatureUnit() | Set the temperature unit. (Celsius, Fahrenheit) | .TemperatureUnit(openmeteogo.Celsius) |
| WindspeedUnit() | Set the wind speed unit. (KMH, MPH, etc.) | .WindspeedUnit(openmeteogo.MPH) |
| PrecipitationUnit() | Set the precipitation unit. (MM, IN) | .PrecipitationUnit(openmeteogo.IN) |
| Timezone() | Set the timezone for results. Takes a `*time.Location`. | .Timezone(time.UTC) |
| PastDays() | Request N number of past days of data. | .PastDays(7) |
| ForecastDays() | Request N number of forecast days. | .ForecastDays(3) |
| Start() | Set a start date for historical queries. | .Start(time.Now()) |
| End() | Set an end date for historical queries. | .End(time.Now()) |
| Seasonal() | Enable Seasonal API. | .Seasonal(true) |
| Marine() | Enable Marine API. | .Marine(true) |
| Models() | Set specific weather models. Accepted on every endpoint, and does not by itself select the Seasonal API. | .Models([]string{"gfs_seamless"}) |
| CurrentMetrics() | Select which current metrics to fetch. | .CurrentMetrics(\&CurrentMetrics{...}) |
| DailyMetrics() | Select which daily metrics to fetch. | .DailyMetrics(\&DailyMetrics{...}) |
| HourlyMetrics() | Select which hourly metrics to fetch. | .HourlyMetrics(\&HourlyMetrics{...}) |
| Marine() | Use the marine endpoint (waves, swell, currents). | .Marine(true) |

### **Available Metrics**

You can select exactly which weather variables you want by populating the CurrentMetrics, DailyMetrics, or HourlyMetrics with a slice of strings, that map directly to the [JSON parameters in the Open-Meteo documentation](https://open-meteo.com/en/docs).


You can request any combination of the following metrics.

### **Current Metrics**

Temperature2m, RelativeHumidity2m, IsDay, ApparentTemperature, Precipitation,
Rain, Showers, Snowfall, WeatherCode, CloudCover, PressureMsl, SurfacePressure,
WindSpeed10m, WindDirection10m, WindGusts10m

### **Hourly Metrics**

Temperature2m, RelativeHumidity2m, DewPoint2m, ApparentTemperature,
PrecipitationProbability, Precipitation, Rain, Showers, Snowfall, SnowDepth,
WeatherCode, PressureMsl, SurfacePressure, CloudCover, CloudCoverLow,
CloudCoverMid, CloudCoverHigh, Evapotranspiration, Visibility,
Et0FaoEvapotranspiration, VapourPressureDeficit, WindSpeed10m, WindSpeed80m,
WindSpeed120m, WindSpeed180m, WindDirection10m, WindDirection80m,
WindDirection120m, WindDirection180m, WindGusts10m, Temperature80m,
Temperature120m, Temperature180m, SoilTemperature0cm, SoilTemperature6cm,
SoilTemperature18cm, SoilTemperature54cm, SoilMoisture0To1cm,
SoilMoisture1To3cm, SoilMoisture9To27cm, SoilMoisture3To9cm,
WaveHeight, WaveDirection, WavePeriod, WavePeakPeriod, WindWaveHeight,
WindWaveDirection, WindWavePeriod, WindWavePeakPeriod, SwellWaveHeight,
SwellWaveDirection, SwellWavePeriod, SwellWavePeakPeriod,
SecondarySwellWaveHeight, SecondarySwellWaveDirection, SecondarySwellWavePeriod,
TertiarySwellWaveHeight, TertiarySwellWaveDirection, TertiarySwellWavePeriod,
SeaLevelHeight, SeaSurfaceTemperature, OceanCurrentVelocity,
OceanCurrentDirection

### **Daily Metrics**

WeatherCode, Temperature2mMax, Temperature2mMin, ApparentTemperatureMax,
ApparentTemperatureMin, Sunrise, Sunset, SunshineDuration, DaylightDuration,
UvIndexMax, UvIndexClearSkyMax, RainSum, ShowersSum, SnowfallSum,
PrecipitationSum, PrecipitationHours, PrecipitationProbabilityMax,
WindSpeed10mMax, WindGusts10mMax, WindDirection10mDominant,
ShortwaveRadiationSum, Et0FaoEvapotranspiration, WaveHeightMax,
WaveDirectionDominant, WavePeriodMax, WindWaveHeightMax,
WindWaveDirectionDominant, WindWavePeriodMax, WindWavePeakPeriodMax,
SwellWaveHeightMax, SwellWaveDirectionDominant, SwellWavePeriodMax,
SwellWavePeakPeriodMax

### **Weekly & Monthly Metrics (Seasonal)**

Temperature2mMean, Temperature2mAnomaly, Temperature2mMaxMean,
Temperature2mMinMean, DewPoint2mMean, PrecipitationMean, PrecipitationAnomaly,
PressureMslMean, PressureMslAnomaly, SoilMoisture0To10cmMean,
SoilMoisture0To10cmAnomaly

### **Marine Metrics**

Only valid with `Marine(true)`. Every one of these decodes to `*float64` /
`[]*float64`.

**Hourly and current** (`marine_hourly`, `marine_current`):

WaveHeight, WaveDirection, WavePeriod, WavePeakPeriod, WindWaveHeight,
WindWaveDirection, WindWavePeriod, WindWavePeakPeriod, SwellWaveHeight,
SwellWaveDirection, SwellWavePeriod, SwellWavePeakPeriod,
SecondarySwellWaveHeight, SecondarySwellWaveDirection, SecondarySwellWavePeriod,
TertiarySwellWaveHeight, TertiarySwellWaveDirection, TertiarySwellWavePeriod,
SeaLevelHeight, SeaSurfaceTemperature, OceanCurrentVelocity,
OceanCurrentDirection

(`WavePeakPeriod` and the secondary/tertiary swell metrics are hourly only.)

**Daily** (`marine_daily`):

WaveHeightMax, WaveDirectionDominant, WavePeriodMax, WindWaveHeightMax,
WindWaveDirectionDominant, WindWavePeriodMax, WindWavePeakPeriodMax,
SwellWaveHeightMax, SwellWaveDirectionDominant, SwellWavePeriodMax,
SwellWavePeakPeriodMax

## **Available Units**

You can specify the units for the following measurements:

* **Temperature**: openmeteogo.Celsius (default), openmeteogo.Fahrenheit  
* **Wind Speed**: openmeteogo.KMH (default), openmeteogo.MS, openmeteogo.MPH, openmeteogo.KN  
* **Precipitation**: openmeteogo.MM (default), openmeteogo.IN




Example:

```go
// Request only current temperature and wind speed
	opts := openmeteogo.NewOptionsBuilder().
		Latitude(37.77).
		Longitude(-122.41).
		TemperatureUnit(openmeteogo.Fahrenheit).
		CurrentMetrics(openmeteogo.Metrics{
			openmeteogo.Temperature2m,
			openmeteogo.WindSpeed10m,
		}).
		Build()
```

This is not an officially supported Google product. This project is not
eligible for the [Google Open Source Software Vulnerability Rewards
Program](https://bughunters.google.com/open-source-security).