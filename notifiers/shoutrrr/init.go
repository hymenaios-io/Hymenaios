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

package shoutrrr

import (
	"strings"

	shoutrrr_types "github.com/containrrr/shoutrrr/pkg/types"
	"github.com/release-argus/Argus/utils"
	metrics "github.com/release-argus/Argus/web/metrics"
)

// Init the Slice metrics.
func (s *Slice) Init(
	log *utils.JLog,
	serviceID *string,
	mains *Slice,
	defaults *Slice,
	hardDefaults *Slice,
) {
	jLog = log
	if s == nil {
		return
	}
	if mains == nil {
		mains = &Slice{}
	}

	for key := range *s {
		id := key
		if (*s)[key] == nil {
			(*s)[key] = &Shoutrrr{}
		}
		(*s)[key].ID = &id
		// Get Type from this or the associated Main
		notifyType := utils.GetFirstNonDefault(
			(*s)[key].Type,
			utils.DefaultIfNil((*mains)[key]).Type,
		)

		// Ensure defaults aren't nil
		if (*defaults)[notifyType] == nil {
			(*defaults)[notifyType] = &Shoutrrr{}
		}
		if (*hardDefaults)[notifyType] == nil {
			(*hardDefaults)[notifyType] = &Shoutrrr{}
		}

		(*s)[key].Init(serviceID, (*mains)[key], (*defaults)[notifyType], (*hardDefaults)[notifyType])
	}
}

// Init the Shoutrrr metrics and hand out the defaults.
func (s *Shoutrrr) Init(
	serviceID *string,
	main *Shoutrrr,
	defaults *Shoutrrr,
	hardDefaults *Shoutrrr,
) {
	s.initMetrics(serviceID)
	s.InitMaps()

	if s == nil {
		s = &Shoutrrr{}
	}

	// Give the matching main
	(*s).Main = main
	if main == nil && utils.DefaultIfNil(serviceID) != "" {
		s.Main = &Shoutrrr{}
	}
	s.Main.InitMaps()

	// Give Defaults
	(*s).Defaults = defaults
	s.Defaults.InitMaps()

	// Give Hard Defaults
	(*s).HardDefaults = hardDefaults
	s.HardDefaults.InitMaps()
}

// initOptions mapping, converting all keys to lowercase.
func (s *Shoutrrr) initOptions() {
	Options := make(map[string]string)
	if s.Options != nil {
		for i := range *s.Options {
			Options[strings.ToLower(i)] = (*s.Options)[i]
		}
	}
	s.Options = &Options
}

// initURLFields mapping, converting all keys to lowercase.
func (s *Shoutrrr) initURLFields() {
	URLFields := make(map[string]string)
	if s.URLFields != nil {
		for i := range *s.URLFields {
			URLFields[strings.ToLower(i)] = (*s.URLFields)[i]
		}
	}
	s.URLFields = &URLFields
}

// initParams mapping, converting all keys to lowercase.
func (s *Shoutrrr) initParams() {
	params := make(shoutrrr_types.Params)
	if s.Params != nil {
		for i := range *s.Params {
			params[strings.ToLower(i)] = (*s.Params)[i]
		}
	}
	s.Params = &params
}

// InitMaps will initialise all maps, converting all keys to lowercase.
func (s *Shoutrrr) InitMaps() {
	if s == nil {
		return
	}
	s.initOptions()
	s.initURLFields()
	s.initParams()
}

// initMetrics, giving them all a starting value.
func (s *Shoutrrr) initMetrics(serviceID *string) {
	if s.Type == "" {
		return
	}

	// ############
	// # Counters #
	// ############
	metrics.InitPrometheusCounterActions(metrics.NotifyMetric, *(*s).ID, *serviceID, s.Type, "SUCCESS")
	metrics.InitPrometheusCounterActions(metrics.NotifyMetric, *(*s).ID, *serviceID, s.Type, "FAIL")
}
