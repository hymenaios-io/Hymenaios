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

package deployed_version

import (
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/release-argus/Argus/service/options"
	"github.com/release-argus/Argus/utils"
)

func TestPrint(t *testing.T) {
	// GIVEN a Lookup
	allowInvalidCerts := false
	tests := map[string]struct {
		lookup    *Lookup
		headers   []Header
		basicAuth *BasicAuth
		options   options.Options
		lines     int
	}{
		"nil lookup": {lines: 0, lookup: nil},
		"lookup with no basicAuth/headers/options": {lines: 5, lookup: &Lookup{
			URL: "https://release-argus.io", AllowInvalidCerts: &allowInvalidCerts, Regex: "[0-9]+", JSON: "version"}},
		"lookup with basicAuth and no headers/options": {lines: 8, basicAuth: &BasicAuth{Username: "foo", Password: "bar"}, lookup: &Lookup{
			URL: "https://release-argus.io", AllowInvalidCerts: &allowInvalidCerts, Regex: "[0-9]+", JSON: "version"}},
		"lookup with headers and no basicAuth/options": {lines: 10, headers: []Header{{Key: "a", Value: "b"}, {Key: "b", Value: "a"}}, lookup: &Lookup{
			URL: "https://release-argus.io", AllowInvalidCerts: &allowInvalidCerts, Regex: "[0-9]+", JSON: "version"}},
		"lookup with no basicAuth/headers": {lines: 5, options: options.Options{Interval: "10m"}, lookup: &Lookup{
			URL: "https://release-argus.io", AllowInvalidCerts: &allowInvalidCerts, Regex: "[0-9]+", JSON: "version"}},
		"lookup with basicAuth and headers": {lines: 13, basicAuth: &BasicAuth{Username: "foo", Password: "bar"},
			options: options.Options{Interval: "10m"}, headers: []Header{{Key: "a", Value: "b"}, {Key: "b", Value: "a"}}, lookup: &Lookup{
				URL: "https://release-argus.io", AllowInvalidCerts: &allowInvalidCerts, Regex: "[0-9]+", JSON: "version"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			stdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			if tc.lookup != nil {
				tc.lookup.Headers = tc.headers
				tc.lookup.BasicAuth = tc.basicAuth
				tc.lookup.Options = &tc.options
			}

			// WHEN Print is called
			tc.lookup.Print("")

			// THEN it prints the expected number of lines
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = stdout
			got := strings.Count(string(out), "\n")
			if got != tc.lines {
				t.Errorf("Print should have given %d lines, but gave %d\n%s",
					tc.lines, got, out)
			}
		})
	}
}

func TestCheckValues(t *testing.T) {
	// GIVEN a Lookup
	tests := map[string]struct {
		url        string
		regex      string
		defaults   *Lookup
		errRegex   string
		nilService bool
	}{
		"nil service":                            {errRegex: `^$`, nilService: true},
		"valid service":                          {errRegex: `^$`, url: "https://example.com", regex: "[0-9.]+", defaults: &Lookup{}},
		"no url":                                 {errRegex: `url: <missing>`, url: "", defaults: &Lookup{}},
		"invalid regex":                          {errRegex: `regex: .* <invalid>`, regex: "[0-", defaults: &Lookup{}},
		"all errs":                               {errRegex: `url: <missing>`, url: "", regex: "[0-", defaults: &Lookup{}},
		"no url doesnt fail for Lookup Defaults": {errRegex: `^$`, url: "", defaults: nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			lookup := &Lookup{}
			*lookup = testDeployedVersion()
			lookup.URL = tc.url
			lookup.Regex = tc.regex
			lookup.Defaults = nil
			if tc.defaults != nil {
				lookup.Defaults = tc.defaults
			}
			if tc.nilService {
				lookup = nil
			}

			// WHEN CheckValues is called
			err := lookup.CheckValues("")

			// THEN it err's when expected
			e := utils.ErrorToString(err)
			re := regexp.MustCompile(tc.errRegex)
			match := re.MatchString(e)
			if !match {
				t.Fatalf("want match for %q\nnot: %q",
					tc.errRegex, e)
			}
		})
	}
}