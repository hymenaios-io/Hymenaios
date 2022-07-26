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

//go:build testing

package latest_version

import (
	db_types "github.com/release-argus/Argus/db/types"
	"github.com/release-argus/Argus/service/latest_version/filters"
	"github.com/release-argus/Argus/service/options"
	service_status "github.com/release-argus/Argus/service/status"
)

func boolPtr(val bool) *bool {
	return &val
}
func stringPtr(val string) *string {
	return &val
}

func testLookupGitHub() Lookup {
	var (
		announceChannel chan []byte           = make(chan []byte, 2)
		saveChannel     chan bool             = make(chan bool, 5)
		databaseChannel chan db_types.Message = make(chan db_types.Message, 5)
	)
	lookup := Lookup{
		Type:              "github",
		URL:               "release-argus/Argus",
		AccessToken:       stringPtr("TOKEN"),
		AllowInvalidCerts: boolPtr(true),
		UsePreRelease:     boolPtr(true),
		Require:           &filters.Require{},
		Options: &options.Options{
			SemanticVersioning: boolPtr(true),
			Defaults:           &options.Options{},
			HardDefaults:       &options.Options{},
		},
		Status: &service_status.Status{
			ServiceID:       stringPtr("test"),
			AnnounceChannel: &announceChannel,
			DatabaseChannel: &databaseChannel,
			SaveChannel:     &saveChannel,
		},
		Defaults:     &Lookup{},
		HardDefaults: &Lookup{},
	}
	lookup.Require.Status = lookup.Status
	return lookup
}

func testLookupURL() Lookup {
	var (
		announceChannel chan []byte           = make(chan []byte, 2)
		saveChannel     chan bool             = make(chan bool, 5)
		databaseChannel chan db_types.Message = make(chan db_types.Message, 5)
	)
	lookup := Lookup{
		Type:              "url",
		URL:               "https://release-argus.io/docs/config/service",
		AllowInvalidCerts: boolPtr(true),
		URLCommands:       filters.URLCommandSlice{{Type: "regex", Regex: stringPtr("[0-9.]+test")}},
		Require:           &filters.Require{},
		Options: &options.Options{
			SemanticVersioning: boolPtr(true),
			Defaults:           &options.Options{},
			HardDefaults:       &options.Options{},
		},
		Status: &service_status.Status{
			ServiceID:       stringPtr("test"),
			AnnounceChannel: &announceChannel,
			DatabaseChannel: &databaseChannel,
			SaveChannel:     &saveChannel,
		},
		Defaults:     &Lookup{},
		HardDefaults: &Lookup{},
	}
	lookup.Require.Status = lookup.Status
	return lookup
}