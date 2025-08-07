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

// TemperatureUnit defines the unit for temperature values.
type TemperatureUnit string

const (
	// Celsius is the default temperature unit.
	Celsius TemperatureUnit = "celsius"
	// Fahrenheit temperature unit.
	Fahrenheit TemperatureUnit = "fahrenheit"
)

// WindSpeedUnit defines the unit for wind speed values.
type WindSpeedUnit string

const (
	// KMH is the default wind speed unit (kilometers per hour).
	KMH WindSpeedUnit = "kmh"
	// MS is meters per second.
	MS WindSpeedUnit = "ms"
	// MPH is miles per hour.
	MPH WindSpeedUnit = "mph"
	// KN is knots.
	KN WindSpeedUnit = "kn"
)

// PrecipitationUnit defines the unit for precipitation values.
type PrecipitationUnit string

const (
	// MM is the default precipitation unit (millimeters).
	MM PrecipitationUnit = "mm"
	// IN is inches.
	IN PrecipitationUnit = "in"
)
