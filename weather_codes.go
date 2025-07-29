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

// WeatherCodeMap maps WMO weather codes to their descriptions.
// Based on WMO Code 4561.
var WeatherCodeMap = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Fog",
	48: "Depositing rime fog",
	51: "Drizzle: Light intensity",
	53: "Drizzle: Moderate intensity",
	55: "Drizzle: Dense intensity",
	56: "Freezing Drizzle: Light intensity",
	57: "Freezing Drizzle: Dense intensity",
	61: "Rain: Slight intensity",
	63: "Rain: Moderate intensity",
	65: "Rain: Heavy intensity",
	66: "Freezing Rain: Light intensity",
	67: "Freezing Rain: Heavy intensity",
	71: "Snow fall: Slight intensity",
	73: "Snow fall: Moderate intensity",
	75: "Snow fall: Heavy intensity",
	77: "Snow grains",
	80: "Rain showers: Slight",
	81: "Rain showers: Moderate",
	82: "Rain showers: Violent",
	85: "Snow showers: Slight",
	86: "Snow showers: Heavy",
	95: "Thunderstorm: Slight or moderate",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

// DescribeCode returns the string description for a given weather code.
// It returns "Unknown code" if the code is not found.
func DescribeCode(code int) string {
	if desc, ok := WeatherCodeMap[code]; ok {
		return desc
	}
	return "Unknown code"
}
