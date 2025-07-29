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

package main

import (
	"fmt"
	"log"
	"time"

	openmeteogo "github.com/tpryan/openmeteogo"
)

func main() {
	// Create a new client
	c := openmeteogo.NewClient()

	// Set the options for the weather data we want.
	// Let's get the current temperature and weather code for San Francisco,
	// as well as the forecast for today.
	opts := openmeteogo.NewOptionsBuilder().
		Latitude(37.7749).Longitude(-122.4194).
		TemperatureUnit(openmeteogo.Fahrenheit).
		ForcastDays(1).
		CurrentMetrics(&openmeteogo.CurrentMetrics{
			Temperature2m: true,
			WeatherCode:   true,
		}).
		DailyMetrics(&openmeteogo.DailyMetrics{
			Temperature2mMax: true,
			Temperature2mMin: true,
			WeatherCode:      true,
		}).
		Build()

	// Make the API call
	weather, err := c.Get(opts)
	if err != nil {
		log.Fatalf("Failed to get weather data: %v", err)
	}

	// Print the current weather
	fmt.Printf("Current Weather in San Francisco (Lat: %.2f, Lon: %.2f)\n", weather.Latitude, weather.Longitude)
	fmt.Printf("Time: %s\n", weather.Current.Time)
	fmt.Printf("Temperature: %.2f%s\n", weather.Current.Temperature2m, weather.CurrentUnits.Temperature2m)
	fmt.Printf("Weather: %s\n\n", openmeteogo.DescribeCode(weather.Current.WeatherCode))

	// Print the daily forecast for today
	if len(weather.Daily.Time) > 0 {
		fmt.Printf("Forecast for %s:\n", weather.Daily.Time[0])
		fmt.Printf("  Max Temperature: %.2f%s\n", weather.Daily.Temperature2mMax[0], weather.DailyUnits.Temperature2mMax)
		fmt.Printf("  Min Temperature: %.2f%s\n", weather.Daily.Temperature2mMin[0], weather.DailyUnits.Temperature2mMin)
		fmt.Printf("  Weather: %s\n", openmeteogo.DescribeCode(weather.Daily.WeatherCode[0]))
	}

	pastOpts := openmeteogo.NewOptionsBuilder().
		Latitude(37.7749).Longitude(-122.4194).
		TemperatureUnit(openmeteogo.Fahrenheit).
		DailyMetrics(&openmeteogo.DailyMetrics{
			Temperature2mMax: true,
			Temperature2mMin: true,
			WeatherCode:      true,
		}).
		Start(time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC)).End(time.Date(2025, 7, 26, 0, 0, 0, 0, time.UTC)).
		Build()

	weatherPast, err := c.Get(pastOpts)
	if err != nil {
		log.Fatalf("Failed to get weather data: %v", err)
	}

	// Print the current weather
	fmt.Printf("Past Weather in San Francisco (Lat: %.2f, Lon: %.2f)\n", weatherPast.Latitude, weatherPast.Longitude)

	// Print the daily forecast for today
	if len(weatherPast.Daily.Time) > 0 {
		fmt.Printf("Forecast for %s:\n", weatherPast.Daily.Time[0])
		fmt.Printf("  Max Temperature: %.2f%s\n", weatherPast.Daily.Temperature2mMax[0], weatherPast.DailyUnits.Temperature2mMax)
		fmt.Printf("  Min Temperature: %.2f%s\n", weatherPast.Daily.Temperature2mMin[0], weatherPast.DailyUnits.Temperature2mMin)
		fmt.Printf("  Weather: %s\n", openmeteogo.DescribeCode(weatherPast.Daily.WeatherCode[0]))
	}

}
