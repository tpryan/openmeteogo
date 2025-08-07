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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemperatureUnitString(t *testing.T) {

	tests := map[string]struct {
		input TemperatureUnit
		want  string
	}{
		"celsius": {
			input: Celsius,
			want:  "celsius",
		},
		"fahrenheit": {
			input: Fahrenheit,
			want:  "fahrenheit",
		},
		"unknown": {
			input: "unknown",
			want:  "unknown",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := string(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPrecipitationUnitString(t *testing.T) {

	tests := map[string]struct {
		input PrecipitationUnit
		want  string
	}{
		"in": {
			input: IN,
			want:  "in",
		},
		"mm": {
			input: MM,
			want:  "mm",
		},
		"unknown": {
			input: "unknown",
			want:  "unknown",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := string(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestWindSpeedUnitString(t *testing.T) {
	tests := map[string]struct {
		input WindSpeedUnit
		want  string
	}{
		"kmh": {
			input: KMH,
			want:  "kmh",
		},
		"ms": {
			input: MS,
			want:  "ms",
		},
		"mph": {
			input: MPH,
			want:  "mph",
		},
		"kn": {
			input: KN,
			want:  "kn",
		},
		"unknown": {
			input: "unknown",
			want:  "unknown",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := string(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}

}
