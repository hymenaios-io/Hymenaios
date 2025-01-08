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
	"os"
	"sync"
	"testing"

	github_types "github.com/release-argus/Argus/service/latest_version/types/github/api_type"
	"github.com/release-argus/Argus/test"
)

var emptyListETagTestMutex = sync.Mutex{}

func Test_SetEmptyListETag(t *testing.T) {
	// GIVEN emptyListETag is set to the incorrect value
	emptyListETagTestMutex.Lock()
	t.Cleanup(func() { emptyListETagTestMutex.Unlock() })
	incorrectValue := "foo"
	setEmptyListETag(incorrectValue)

	// WHEN SetEmptyListETag is called
	SetEmptyListETag(os.Getenv("GITHUB_TOKEN"))

	// THEN the emptyListETag is set
	setTo := getEmptyListETag()
	if setTo == incorrectValue {
		t.Errorf("emptyListETag wasn't updated. Got %q, want %q",
			setTo, emptyListETag)
	}
	if setTo != initialEmptyListETag {
		t.Errorf("Empty list ETag has changed from %q to %q",
			initialEmptyListETag, setTo)
	}
}

func Test_setEmptyListETag(t *testing.T) {
	// GIVEN emptyListETag exists
	emptyListETagTestMutex.Lock()
	t.Cleanup(func() { emptyListETagTestMutex.Unlock() })

	// WHEN setEmptyListETag is called
	newValue := "foo"
	setEmptyListETag(newValue)

	// THEN the emptyListETag is set
	if emptyListETag != newValue {
		t.Errorf("setEmptyListETag() = %q, want %q",
			emptyListETag, newValue)
	}
}

func TestGetEmptyListETag(t *testing.T) {
	// GIVEN emptyListETag exists
	emptyListETagTestMutex.Lock()
	t.Cleanup(func() { emptyListETagTestMutex.Unlock() })
	emptyListETagMutex.RLock()
	t.Cleanup(func() { emptyListETagMutex.RUnlock() })

	// WHEN getEmptyListETag is called
	got := getEmptyListETag()

	// THEN the emptyListETag is returned
	if got != emptyListETag {
		t.Errorf("getEmptyListETag() = %q, want %q", got, emptyListETag)
	}
}

func TestNewData(t *testing.T) {
	emptyListETagTestMutex.Lock()
	t.Cleanup(func() { emptyListETagTestMutex.Unlock() })
	startingEmptyListETag := getEmptyListETag()
	// GIVEN a Data is wanted with/without an eTag/releases
	tests := map[string]struct {
		eTag     string
		releases *[]github_types.Release
		want     *Data
	}{
		"no eTag or releases": {
			eTag:     "",
			releases: nil,
			want: &Data{
				eTag:     startingEmptyListETag,
				releases: []github_types.Release{},
			},
		},
		"eTag but no releases": {
			eTag:     "foo",
			releases: nil,
			want: &Data{
				eTag:     "foo",
				releases: []github_types.Release{},
			},
		},
		"no eTag but releases": {
			eTag: "",
			releases: &[]github_types.Release{
				{TagName: "bar"}},
			want: &Data{
				eTag: startingEmptyListETag,
				releases: []github_types.Release{
					{TagName: "bar"}},
			},
		},
		"eTag and releases": {
			eTag: "zing",
			releases: &[]github_types.Release{
				{TagName: "zap"}},
			want: &Data{
				eTag: "zing",
				releases: []github_types.Release{
					{TagName: "zap"}},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN newData is called
			got := newData(tc.eTag, tc.releases)

			// THEN the correct Data is returned
			if got.eTag != tc.want.eTag {
				t.Errorf("eTag: got %q, want %q",
					got.eTag, tc.want.eTag)
			}
			if len(got.releases) != len(tc.want.releases) {
				t.Errorf("releases: got %v, want %v",
					got.releases, tc.want.releases)
			} else {
				for i, release := range got.releases {
					if release.TagName != tc.want.releases[i].TagName {
						t.Errorf("%d: TagName, got %q (%v), want %q (%v)",
							i, got.releases[i].TagName, got.releases, tc.want.releases[i].TagName, tc.want.releases)
					}
				}
			}
		})
	}
}

func TestData_String(t *testing.T) {
	// GIVEN a Data
	tests := map[string]struct {
		githubData *Data
		want       string
	}{
		"nil": {
			githubData: nil,
			want:       ""},
		"empty": {
			githubData: &Data{},
			want:       "{}"},
		"filled": {
			githubData: &Data{
				eTag: "argus",
				releases: []github_types.Release{
					{URL: "https://test.com/1.2.3"},
					{URL: "https://test.com/3.2.1", PreRelease: true},
				}},
			want: `
				{
					"etag": "argus",
					"releases": [
						{"url": "https://test.com/1.2.3", "prerelease": false},
						{"url": "https://test.com/3.2.1", "prerelease": true}
					]
				}`},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc.want = test.TrimJSON(tc.want)

			// WHEN the Data is stringified with String
			got := tc.githubData.String()

			// THEN the result is as expected
			if got != tc.want {
				t.Errorf("got:\n%q\nwant:\n%q",
					got, tc.want)
			}
		})
	}
}

func TestData_TagFallback(t *testing.T) {
	// GIVEN a Data
	gd := &Data{}
	tests := []bool{
		true, false, true, false, true}

	if gd.tagFallback != false {
		t.Fatalf("tagFallback wasn't set to false initially")
	}

	for _, tc := range tests {
		gd.SetTagFallback()

		// WHEN TagFallback is called
		got := gd.TagFallback()

		// THEN the correct value is returned
		if got != tc {
			t.Errorf("got %t, want %t", got, tc)
		}
	}
}

func TestData_ETag(t *testing.T) {
	// GIVEN a Data
	test := &Data{}

	// WHEN ETag is called
	got := test.ETag()

	// THEN the releases are returned
	want := test.eTag
	if got != want {
		t.Errorf("got %q, want %q",
			got, want)
	}

	// WHEN the releases are changed
	newETag := "foo"
	test.SetETag(newETag)

	// THEN the new releases can be fetched
	got = test.ETag()
	want = newETag
	if got != want {
		t.Errorf("got %q, want %q",
			got, want)
	}
}

func TestData_Releases(t *testing.T) {
	// GIVEN a Data
	test := &Data{}

	// WHEN Releases is called
	got := test.Releases()

	// THEN the releases are returned
	want := test.releases
	match := len(got) == len(want)
	if match {
		for i, release := range got {
			if release.String() != want[i].String() {
				match = false
				break
			}
		}
	}
	if !match {
		t.Errorf("got %v, want %v",
			got, want)
	}

	// WHEN the releases are changed
	newReleases := []github_types.Release{
		{TagName: "foo"},
		{TagName: "bar"}}
	test.SetReleases(newReleases)

	// THEN the new releases can be fetched
	got = test.Releases()
	want = newReleases
	match = len(got) == len(want)
	if match {
		for i, release := range got {
			if release.String() != want[i].String() {
				match = false
				break
			}
		}
	}
	if !match {
		t.Errorf("got %v, want %v",
			got, want)
	}
}

func TestData_hasReleases(t *testing.T) {
	// GIVEN a Data that may/may not have releases
	tests := map[string]struct {
		gd   *Data
		want bool
	}{
		"no releases": {
			gd:   &Data{},
			want: false,
		},
		"1 release": {
			gd: &Data{
				releases: []github_types.Release{
					{TagName: "foo"}}},
			want: true,
		},
		"multiple releases": {
			gd: &Data{
				releases: []github_types.Release{
					{TagName: "foo"},
					{TagName: "bar"}}},
			want: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN hasReleases is called on it
			got := tc.gd.hasReleases()

			// THEN the correct value is returned
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestData_Copy(t *testing.T) {
	// GIVEN a Data to copy from
	tests := map[string]struct {
		gd *Data
	}{
		"empty": {
			gd: &Data{},
		},
		"filled": {
			gd: &Data{
				eTag: "foo",
				releases: []github_types.Release{
					{TagName: "bar"}}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN Copy is called
			got := tc.gd.Copy()

			// THEN the correct Data is returned
			if got.eTag != tc.gd.eTag {
				t.Errorf("eTag: got %q, want %q",
					got.eTag, tc.gd.eTag)
			}
			if len(got.releases) != len(tc.gd.releases) {
				t.Errorf("releases: got %v, want %v",
					got.releases, tc.gd.releases)
			} else {
				for i, release := range got.releases {
					if release.TagName != tc.gd.releases[i].TagName {
						t.Errorf("%d: TagName, got %q (%v), want %q (%v)",
							i, got.releases[i].TagName, got.releases, tc.gd.releases[i].TagName, tc.gd.releases)
					}
				}
			}
		})
	}
}

func TestData_CopyFrom(t *testing.T) {
	// GIVEN a fresh Data and a Data to copy from
	tests := map[string]struct {
		fresh *Data
		gd    *Data
	}{
		"empty": {
			gd: &Data{},
		},
		"filled": {
			gd: &Data{
				eTag: "foo",
				releases: []github_types.Release{
					{TagName: "bar"}}},
		},
		"filled with data to overwrite": {
			fresh: &Data{
				eTag: "fizz",
				releases: []github_types.Release{
					{TagName: "bang"}}},
			gd: &Data{
				eTag: "foo",
				releases: []github_types.Release{
					{TagName: "bar"}}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tc.fresh == nil {
				tc.fresh = &Data{}
			}

			// WHEN CopyFrom is called
			tc.fresh.CopyFrom(tc.gd)

			// THEN the correct Data is returned
			if tc.fresh.eTag != tc.gd.eTag {
				t.Errorf("eTag: got %q, want %q",
					tc.fresh.eTag, tc.gd.eTag)
			}
			if len(tc.fresh.releases) != len(tc.gd.releases) {
				t.Errorf("releases: got %v, want %v",
					tc.fresh.releases, tc.gd.releases)
			} else {
				for i, release := range tc.fresh.releases {
					if release.TagName != tc.gd.releases[i].TagName {
						t.Errorf("%d: TagName, got %q (%v), want %q (%v)",
							i, tc.fresh.releases[i].TagName, tc.fresh.releases, tc.gd.releases[i].TagName, tc.gd.releases)
					}
				}
			}
		})
	}
}