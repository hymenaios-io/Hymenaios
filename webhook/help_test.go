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

//go:build unit || integration

package webhook

import (
	"os"
	"strings"
	"testing"

	"github.com/release-argus/Argus/notify/shoutrrr"
	"github.com/release-argus/Argus/service/status"
	"github.com/release-argus/Argus/test"
	"github.com/release-argus/Argus/util"
)

func TestMain(m *testing.M) {
	// initialise jLog
	mainJLog := util.NewJLog("DEBUG", false)
	mainJLog.Testing = true
	LogInit(mainJLog)
	shoutrrr.LogInit(mainJLog)

	// run other tests
	exitCode := m.Run()

	// exit
	os.Exit(exitCode)
}

func testWebHook(failing bool, selfSignedCert bool, customHeaders bool) *WebHook {
	desiredStatusCode := uint16(0)
	whMaxTries := uint8(1)
	webhook := New(
		test.BoolPtr(false),
		nil,
		"0s",
		&desiredStatusCode,
		nil,
		&whMaxTries,
		nil,
		test.StringPtr("12m"),
		"argus",
		test.BoolPtr(false),
		"github",
		"https://valid.release-argus.io/hooks/github-style",
		&Defaults{},
		&Defaults{}, &Defaults{})
	webhook.ID = "test"
	webhook.ServiceStatus = &status.Status{}
	webhook.ServiceStatus.Init(
		0, 0, 1,
		test.StringPtr("testServiceID"),
		nil)
	webhook.Failed = &webhook.ServiceStatus.Fails.WebHook
	serviceName := "testServiceID"
	webURL := "https://example.com"
	webhook.ServiceStatus.Init(
		0, 1, 0,
		&serviceName,
		&webURL)
	if selfSignedCert {
		webhook.URL = strings.Replace(webhook.URL, "valid", "invalid", 1)
	}
	if failing {
		webhook.Secret = "invalid"
	}
	if customHeaders {
		webhook.URL = strings.Replace(webhook.URL, "github-style", "single-header", 1)
		if failing {
			webhook.CustomHeaders = &Headers{
				{Key: "X-Test", Value: "invalid"}}
		} else {
			webhook.CustomHeaders = &Headers{
				{Key: "X-Test", Value: "secret"}}
		}
	}
	webhook.initMetrics()
	return webhook
}

func testDefaults(failing bool, customHeaders bool) *Defaults {
	desiredStatusCode := uint16(0)
	whMaxTries := uint8(1)
	webhook := NewDefaults(
		test.BoolPtr(false),
		nil,
		"0s",
		&desiredStatusCode,
		&whMaxTries,
		"argus",
		test.BoolPtr(false),
		"github",
		"https://valid.release-argus.io/hooks/github-style")
	if failing {
		webhook.Secret = "invalid"
	}
	if customHeaders {
		webhook.URL = strings.Replace(webhook.URL, "github-style", "single-header", 1)
		if failing {
			webhook.CustomHeaders = &Headers{
				{Key: "X-Test", Value: "invalid"}}
		} else {
			webhook.CustomHeaders = &Headers{
				{Key: "X-Test", Value: "secret"}}
		}
	}
	return webhook
}
