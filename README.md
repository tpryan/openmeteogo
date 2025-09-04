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
* Fluent **Options Builder** for easily constructing complex queries.  
* Support for multiple **units** (Celsius/Fahrenheit, km/h/mph, etc.).  
* Select exactly which weather metrics you need.  
* Helper function to get human-readable descriptions from WMO weather codes.  
* Supports both the free and commercial (API key) Open-Meteo endpoints.

## **Installation**

To install the library, use go get:

```bash
go get github.com/tpryan/openmeteogo
```

## **Usage**

Here is a quick example of how to get the current weather and today's forecast
for San Francisco.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tpryan/openmeteogo"
)

func main() {
	// 1. Create a new client
	c := openmeteogo.NewClient()

	// 2. Set the options for the weather data you want.
	// Use the fluent builder to easily configure your request.
	opts := openmeteogo.NewOptionsBuilder().
		Latitude(37.7749).
		Longitude(-122.4194).
		TemperatureUnit(openmeteogo.Fahrenheit).
		ForcastDays(1).
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
	w, err := c.Get(opts)
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

	wp, err := c.Get(pastOpts)
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

## **Options**

The OptionsBuilder provides a simple way to configure your request.

| Method | Description | Example |
| :---- | :---- | :---- |
| Latitude() | Set the geographical latitude. | .Latitude(37.7749) |
| Longitude() | Set the geographical longitude. | .Longitude(-122.4194) |
| TemperatureUnit() | Set the temperature unit. (Celsius, Fahrenheit) | .TemperatureUnit(openmeteogo.Celsius) |
| WindspeedUnit() | Set the wind speed unit. (KMH, MPH, etc.) | .WindspeedUnit(openmeteogo.MPH) |
| PrecipitationUnit() | Set the precipitation unit. (MM, IN) | .PrecipitationUnit(openmeteogo.IN) |
| Timezone() | Set the timezone for results. | .Timezone(\*time.UTC) |
| PastDays() | Request N number of past days of data. | .PastDays(7) |
| ForcastDays() | Request N number of forecast days. | .ForcastDays(3) |
| Start() | Set a start date for historical queries. | .Start(time.Now()) |
| End() | Set an end date for historical queries. | .End(time.Now()) |
| CurrentMetrics() | Select which current metrics to fetch. | .CurrentMetrics(\&CurrentMetrics{...}) |
| DailyMetrics() | Select which daily metrics to fetch. | .DailyMetrics(\&DailyMetrics{...}) |
| HourlyMetrics() | Select which hourly metrics to fetch. | .HourlyMetrics(\&HourlyMetrics{...}) |

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
SoilMoisture1To3cm, SoilMoisture9To27cm, SoilMoisture3To9cm

### **Daily Metrics**

WeatherCode, Temperature2mMax, Temperature2mMin, ApparentTemperatureMax,
ApparentTemperatureMin, Sunrise, Sunset, SunshineDuration, DaylightDuration,
UvIndexMax, UvIndexClearSkyMax, RainSum, ShowersSum, SnowfallSum,
PrecipitationSum, PrecipitationHours, PrecipitationProbabilityMax,
WindSpeed10mMax, WindGusts10mMax, WindDirection10mDominant,
ShortwaveRadiationSum, Et0FaoEvapotranspiration

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