// Copyright [2022] [Argus]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build unit

package util

import (
	"testing"
)

func TestRegexCheck(t *testing.T) {
	// GIVEN a variety of RegEx's to apply to a string
	str := `testing\n"beta-release": "0.1.2-beta"\n"stable-release": "0.1.1"`
	tests := map[string]struct {
		regex string
		match bool
	}{
		"regex match":    {regex: `release": "[0-9.]+"`, match: true},
		"no regex match": {regex: `release": "[0-9.]+-alpha"`, match: false},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN RegexCheck is called
			got := RegexCheck(tc.regex, str)

			// THEN the regex matches when expected
			if got != tc.match {
				t.Errorf("wanted match=%t, not %t\n%q on %q",
					tc.match, got, tc.regex, str)
			}
		})
	}
}

func TestRegexCheckWithParams(t *testing.T) {
	// GIVEN a variety of RegEx's to apply to a string
	str := `testing\n"beta-release": "0.1.2-beta"\n"stable-release": "0.1.1"`
	tests := map[string]struct {
		regex   string
		version string
		match   bool
	}{
		"regex match": {
			regex:   `release": "{{ version }}"`,
			version: "0.1.1",
			match:   true},
		"no regex match": {
			regex:   `release": "{{ version }}"`,
			version: "0.1.2",
			match:   false},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN RegexCheck is called
			got := RegexCheckWithParams(tc.regex, str, tc.version)

			// THEN the regex matches when expected
			if got != tc.match {
				t.Errorf("wanted match=%t, not %t\n%q on %q",
					tc.match, got, tc.regex, str)
			}
		})
	}
}
