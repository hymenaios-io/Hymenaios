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

package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestInitPrometheusCounterVec(t *testing.T) {
	// GIVEN a metric
	tests := map[string]struct {
		metric *prometheus.CounterVec
		args   []string
	}{
		"LatestVersionQueryMetric":   {metric: LatestVersionQueryMetric, args: []string{"SERVICE_ID", "RESULT"}},
		"DeployedVersionQueryMetric": {metric: DeployedVersionQueryMetric, args: []string{"SERVICE_ID", "RESULT"}},
		"CommandMetric":              {metric: CommandMetric, args: []string{"COMMAND_ID", "RESULT", "SERVICE_ID"}},
		"NotifyMetric":               {metric: NotifyMetric, args: []string{"NOTIFY_ID", "RESULT", "SERVICE_ID", "TYPE"}},
		"WebHookMetric":              {metric: WebHookMetric, args: []string{"WEBHOOK_ID", "RESULT", "SERVICE_ID"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := testutil.CollectAndCount(tc.metric)
			want := 0
			if got != want {
				t.Errorf("%s:\nhaven't initialised yet but got %d metrics, expecting %d",
					name, got, want)
			}

			// WHEN it's initialised with InitPrometheusCounterWithIDAndResult
			if len(tc.args) == 2 {
				InitPrometheusCounterWithIDAndResult(tc.metric, tc.args[0], tc.args[1])
			} else if len(tc.args) == 3 {
				InitPrometheusCounterActions(tc.metric, tc.args[0], tc.args[2], "", tc.args[1])
			} else {
				InitPrometheusCounterActions(tc.metric, tc.args[0], tc.args[2], tc.args[3], tc.args[1])
			}
			got = testutil.CollectAndCount(tc.metric)
			want = 1
			if got != want {
				t.Errorf("%s:\nhas been initialised but got %d metrics, expecting %d",
					name, got, want)
			}
			var wantValue float64
			var gotValue float64
			if len(tc.args) == 2 {
				gotValue = testutil.ToFloat64(tc.metric.WithLabelValues(tc.args[0], tc.args[1]))
			} else if len(tc.args) == 3 {
				gotValue = testutil.ToFloat64(tc.metric.WithLabelValues(tc.args[0], tc.args[1], tc.args[2]))
			} else {
				gotValue = testutil.ToFloat64(tc.metric.WithLabelValues(tc.args[0], tc.args[1], tc.args[2], tc.args[3]))
			}
			if gotValue != wantValue {
				t.Errorf("%s:\nhas been initialised but got %f, expecting %f",
					name, gotValue, wantValue)
			}

			// THEN it can be increased
			if len(tc.args) == 2 {
				IncreasePrometheusCounterWithIDAndResult(tc.metric, tc.args[0], tc.args[1])
			} else if len(tc.args) == 3 {
				IncreasePrometheusCounterActions(tc.metric, tc.args[0], tc.args[2], "", tc.args[1])
			} else {
				IncreasePrometheusCounterActions(tc.metric, tc.args[0], tc.args[2], tc.args[3], tc.args[1])
			}
			if len(tc.args) == 2 {
				gotValue = testutil.ToFloat64(tc.metric.WithLabelValues(tc.args[0], tc.args[1]))
			} else if len(tc.args) == 3 {
				gotValue = testutil.ToFloat64(tc.metric.WithLabelValues(tc.args[0], tc.args[1], tc.args[2]))
			} else {
				gotValue = testutil.ToFloat64(tc.metric.WithLabelValues(tc.args[0], tc.args[1], tc.args[2], tc.args[3]))
			}
			wantValue++
			if gotValue != wantValue {
				t.Errorf("%s:\nhas been changed but got %f, expecting %f",
					name, gotValue, wantValue)
			}
		})
	}
}

func TestPrometheusGaugeVec(t *testing.T) {
	// GIVEN a metric
	tests := map[string]struct {
		metric     *prometheus.GaugeVec
		args       []string
		isGaugeVec bool
		value      float64
	}{
		"LatestVersionQueryLiveness":   {metric: LatestVersionQueryLiveness, args: []string{"SERVICE_ID"}},
		"DeployedVersionQueryLiveness": {metric: DeployedVersionQueryLiveness, args: []string{"SERVICE_ID"}},
		"AckWaiting":                   {metric: AckWaiting, args: []string{"SERVICE_ID"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := testutil.CollectAndCount(tc.metric)
			want := 0
			if got != want {
				t.Errorf("%s:\nhaven't initialised yet but got %d metrics, expecting %d",
					name, got, want)
			}

			// WHEN it's initialised with SetPrometheusGaugeWithID
			wantValue := float64(3)
			SetPrometheusGaugeWithID(tc.metric, tc.args[0], wantValue)
			got = testutil.CollectAndCount(tc.metric)
			want = 1
			if got != want {
				t.Errorf("%s:\nhas been initialised but got %d metrics, expecting %d",
					name, got, want)
			}
			gotValue := testutil.ToFloat64(tc.metric.WithLabelValues(tc.args[0]))
			if gotValue != wantValue {
				t.Errorf("%s:\nhas been initialised but got %f, expecting %f",
					name, gotValue, wantValue)
			}

			// THEN changes can be noticed
			wantValue = float64(0)
			SetPrometheusGaugeWithID(tc.metric, tc.args[0], wantValue)
			gotValue = testutil.ToFloat64(tc.metric.WithLabelValues(tc.args[0]))
			if gotValue != wantValue {
				t.Errorf("%s:\nhas been changed but got %f, expecting %f",
					name, gotValue, wantValue)
			}
		})
	}
}
