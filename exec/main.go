package main

import (
	"fmt"
	"log"

	"github.com/tpryan/openmateo"
)

func main() {
	// Create a new client
	c := openmateo.New()

	// Set the options for the weather data we want.
	// Let's get the current temperature and weather code for New York City,
	// as well as the forecast for today.
	opts := &openmateo.Options{
		Latitude:    40.71, // New York City
		Longitude:   -74.01,
		ForcastDays: 1, // We want today's forecast
		CurrentMetrics: &openmateo.CurrentMetrics{
			Temperature2m: true,
			WeatherCode:   true,
		},
		DailyMetrics: &openmateo.DailyMetrics{
			Temperature2mMax: true,
			Temperature2mMin: true,
			WeatherCode:      true,
		},
	}

	// Make the API call
	weather, err := c.Get(opts)
	if err != nil {
		log.Fatalf("Failed to get weather data: %v", err)
	}

	// Print the current weather
	fmt.Printf("Current Weather in New York City (Lat: %.2f, Lon: %.2f)\n", weather.Latitude, weather.Longitude)
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
}
