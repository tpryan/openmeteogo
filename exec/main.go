package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tpryan/openmateo"
)

func main() {
	// Create a new client
	c := openmateo.NewClient()

	// Set the options for the weather data we want.
	// Let's get the current temperature and weather code for San Francisco,
	// as well as the forecast for today.
	opts := openmateo.NewOptionsBuilder().
		Latitude(37.7749).Longitude(-122.4194).
		TemperatureUnit(openmateo.Fahrenheit).
		ForcastDays(1).
		CurrentMetrics(&openmateo.CurrentMetrics{
			Temperature2m: true,
			WeatherCode:   true,
		}).
		DailyMetrics(&openmateo.DailyMetrics{
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
	fmt.Printf("Weather Code: %d\n\n", weather.Current.WeatherCode)

	// Print the daily forecast for today
	if len(weather.Daily.Time) > 0 {
		fmt.Printf("Forecast for %s:\n", weather.Daily.Time[0])
		fmt.Printf("  Max Temperature: %.2f%s\n", weather.Daily.Temperature2mMax[0], weather.DailyUnits.Temperature2mMax)
		fmt.Printf("  Min Temperature: %.2f%s\n", weather.Daily.Temperature2mMin[0], weather.DailyUnits.Temperature2mMin)
		fmt.Printf("  Weather Code: %d\n", weather.Daily.WeatherCode[0])
	}

	pastOpts := openmateo.NewOptionsBuilder().
		Latitude(37.7749).Longitude(-122.4194).
		TemperatureUnit(openmateo.Fahrenheit).
		DailyMetrics(&openmateo.DailyMetrics{
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
		fmt.Printf("  Weather Code: %d\n", weatherPast.Daily.WeatherCode[0])
	}

}
