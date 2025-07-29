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

package openmateo

// TemperatureUnit defines the unit for temperature values.
type TemperatureUnit int64

const (
	// Celsius is the default temperature unit.
	Celsius TemperatureUnit = iota
	// Fahrenheit temperature unit.
	Fahrenheit
)

// String returns the string representation of the TemperatureUnit for the API.
func (t TemperatureUnit) String() string {
	switch t {
	case Celsius:
		return "celsius"
	case Fahrenheit:
		return "fahrenheit"
	}
	return "unknown"
}

// WindSpeedUnit defines the unit for wind speed values.
type WindSpeedUnit int64

const (
	// KMH is the default wind speed unit (kilometers per hour).
	KMH WindSpeedUnit = iota
	// MS is meters per second.
	MS
	// MPH is miles per hour.
	MPH
	// KN is knots.
	KN
)

// String returns the string representation of the WindSpeedUnit for the API.
func (w WindSpeedUnit) String() string {
	switch w {
	case KMH:
		return "kmh"
	case MS:
		return "ms"
	case MPH:
		return "mph"
	case KN:
		return "kn"
	}
	return "unknown"
}

// PrecipitationUnit defines the unit for precipitation values.
type PrecipitationUnit int64

const (
	// MM is the default precipitation unit (millimeters).
	MM PrecipitationUnit = iota
	// IN is inches.
	IN
)

// String returns the string representation of the PrecipitationUnit for the API.
func (p PrecipitationUnit) String() string {
	switch p {
	case MM:
		return "mm"
	case IN:
		return "in"
	}
	return "unknown"
}
