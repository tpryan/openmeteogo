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

func TestGetWeatherCodeDescription(t *testing.T) {
	tests := map[string]struct {
		code int
		want string
	}{
		"Clear sky": {
			code: 0,
			want: "Clear sky",
		},
		"Unknown code": {
			code: 100,
			want: "Unknown code",
		},
		"Negative code": {
			code: -1,
			want: "Unknown code",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := DescribeCode(tc.code)
			assert.Equal(t, tc.want, got)
		})
	}
}
