// Copyright [2024] [Argus]
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

// Package github provides a github-based lookup type.
package github

import (
	"testing"
	"time"

	"github.com/release-argus/Argus/service/latest_version/filter"
	"github.com/release-argus/Argus/service/latest_version/types/base"
	github_types "github.com/release-argus/Argus/service/latest_version/types/github/api_type"
	"github.com/release-argus/Argus/test"
	"github.com/release-argus/Argus/util"
	"gopkg.in/yaml.v3"
)

func checkDockerToken(t *testing.T, gotQueryToken, wantQueryToken string, gotValidUntil, wantValidUntil time.Time, message string) {
	if gotQueryToken != wantQueryToken {
		t.Errorf("Require.Docker.queryToken %s\nwant: %q\ngot:  %q",
			message, wantQueryToken, gotQueryToken)
	}
	if gotValidUntil != wantValidUntil {
		t.Errorf("Require.Docker.validUntil %s\nwant: %q\ngot:  %q",
			message, wantValidUntil, gotValidUntil)
	}
}

func TestLookup_Inherit(t *testing.T) {
	testData := newData(
		"etag",
		&[]github_types.Release{
			{URL: "foo"},
			{URL: "bar"}})
	testData.SetTagFallback()
	testRequire := &filter.Require{
		Docker: filter.NewDockerCheck(
			"ghcr",
			"release-argus/argus", "{{ version }}",
			"ghcr-username", "ghcr-token",
			"ghcr-query-token", time.Now(),
			nil)}

	// GIVEN a Lookup and a Lookup to inherit from
	tests := map[string]struct {
		typeChanged    bool
		overrides      string
		inheritData    bool
		inheritRequire bool
	}{
		"don't inherit Data as Type changed": {
			typeChanged: true,
			overrides: test.TrimYAML(`
				type: something-else
				url: something-else
			`),
			inheritData: false,
		},
		"don't inherit Data as URL changed": {
			overrides: test.TrimYAML(`
				url: something-else
			`),
			inheritData: false,
		},
		"inherit Data": {
			overrides: test.TrimYAML(`
				require:
					docker:
						type: hub
			`),
			inheritData: true,
		},
		"inherit Require, not Data": {
			overrides: test.TrimYAML(`
				url: something-else
			`),
			inheritRequire: true,
		},
		"don't inherit Require as Docker changed": {
			overrides: test.TrimYAML(`
				require:
					docker:
						type: something-else
			`),
			inheritData:    true,
			inheritRequire: false,
		},
		"inherit all": {
			inheritData:    true,
			inheritRequire: true,
		},
	}

	otherLookupTest := &struct {
		base.Lookup `yaml:",inline"`
		Data        *Data `yaml:"github_data"`
	}{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			toLookup := testLookup(false)
			testRequireCopy := *testRequire
			toLookup.Require = &testRequireCopy
			toLookup.Require.Docker.SetQueryToken(
				"",
				"", time.Time{})
			toLookup.data.eTag = ""
			var fromLookup base.Interface
			if !tc.typeChanged {
				ghFromLookup := testLookup(false)
				ghFromLookup.data.CopyFrom(testData)
				fromLookup = ghFromLookup
			} else {
				otherLookup := *otherLookupTest
				fromLookup = &otherLookup
			}
			err := yaml.Unmarshal([]byte(tc.overrides), fromLookup)
			if err != nil {
				t.Fatalf("error unmarshalling overrides: %v",
					err)
			}

			// WHEN we call Inherit
			toLookup.Inherit(fromLookup)

			// THEN the Data is copied when expected
			if tc.inheritData {
				if toLookup.data.ETag() != testData.ETag() {
					t.Errorf("Data.ETag not copied over\nwant: %q\ngot;  %q",
						testData.ETag(), toLookup.data.ETag())
				}
				if util.ToYAMLString(toLookup.data.Releases(), "") != util.ToYAMLString(testData.releases, "") {
					t.Errorf("Data.Releases not copied over\nwant: %q\ngot;  %q",
						util.ToYAMLString(testData.releases, ""), util.ToYAMLString(toLookup.data.Releases(), ""))
				}
			} else if toLookup.data.ETag() != "" {
				t.Errorf("Data shouldn't have changed\nwant: %q\ngot;  %q",
					"", toLookup.data.ETag())
			}
			// AND the Require is copied when expected
			if tc.inheritRequire {
				if toLookup.Require == nil {
					t.Errorf("Require not copied over, it's nil")
				} else if toLookup.Require.Docker == nil {
					t.Errorf("Require.Docker not copied over, it's nil")
				} else {
					gotQueryToken, gotValidUntil := toLookup.Require.Docker.CopyQueryToken()
					wantQueryToken, wantValidUntil := testRequire.Docker.CopyQueryToken()
					checkDockerToken(t,
						gotQueryToken, wantQueryToken,
						gotValidUntil, wantValidUntil,
						"not copied")
				}
			} else if toLookup.Require != nil && toLookup.Require.Docker != nil {
				gotQueryToken, gotValidUntil := toLookup.Require.Docker.CopyQueryToken()
				wantQueryToken, wantValidUntil := "", time.Time{}
				checkDockerToken(t,
					gotQueryToken, wantQueryToken,
					gotValidUntil, wantValidUntil,
					"should not be copied")
			}
		})
	}
}