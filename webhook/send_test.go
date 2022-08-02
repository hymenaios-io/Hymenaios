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

package webhook

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/release-argus/Argus/notifiers/shoutrrr"
	"github.com/release-argus/Argus/utils"
)

func TestTry(t *testing.T) {
	// GIVEN a WebHook
	testLogging("WARN")
	tests := map[string]struct {
		url               *string
		allowInvalidCerts bool
		selfSignedCert    bool
		wouldFail         bool
		errRegex          string
		desiredStatusCode int
	}{
		"invalid url": {url: stringPtr("invalid://	test"), errRegex: "failed to get .?http.request"},
		"fail due to invalid secret":              {wouldFail: true, errRegex: "WebHook gave [0-9]+, not "},
		"fail due to invalid cert":                {selfSignedCert: true, errRegex: "certificate signed by unknown authority"},
		"pass with invalid certs allowed":         {selfSignedCert: true, errRegex: "^$", allowInvalidCerts: true},
		"pass with valid certs":                   {errRegex: "^$", allowInvalidCerts: true},
		"fail by not getting desired status code": {desiredStatusCode: 1, errRegex: "WebHook gave [0-9]+, not ", allowInvalidCerts: true},
		"pass by getting desired status code":     {wouldFail: true, desiredStatusCode: 500, errRegex: "^$", allowInvalidCerts: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			webhook := testWebHook(false, true, tc.selfSignedCert)
			if tc.wouldFail {
				webhook = testWebHook(true, true, tc.selfSignedCert)
			}
			if tc.url != nil {
				webhook.URL = *tc.url
			}
			webhook.AllowInvalidCerts = &tc.allowInvalidCerts
			webhook.DesiredStatusCode = &tc.desiredStatusCode

			// WHEN try is called with it
			err := webhook.try(utils.LogFrom{})

			// THEN any err is expected
			e := utils.ErrorToString(err)
			re := regexp.MustCompile(tc.errRegex)
			match := re.MatchString(e)
			if !match {
				t.Errorf("%s:\nwant match for %q\nnot: %q",
					name, tc.errRegex, e)
			}
		})
	}
}

func TestWebHookSend(t *testing.T) {
	// GIVEN a WebHook
	testLogging("INFO")
	tests := map[string]struct {
		wouldFail   bool
		useDelay    bool
		delay       string
		stdoutRegex string
		tries       int
		silentFails bool
		notifiers   shoutrrr.Slice
	}{
		"successful webhook":                           {stdoutRegex: "WebHook received"},
		"does use delay webhook":                       {stdoutRegex: "WebHook received"},
		"failing webhook":                              {wouldFail: true, stdoutRegex: `failed \d times to send`},
		"tries multiple times":                         {wouldFail: true, tries: 2, stdoutRegex: `(WebHook gave 500.*){2}WebHook received`},
		"does try notifiers on fail":                   {wouldFail: true, stdoutRegex: `WebHook gave 500.*invalid gotify token`, notifiers: shoutrrr.Slice{"fail": testNotifier(true, false)}},
		"doesn't try notifiers on fail if silentFails": {wouldFail: true, silentFails: true, stdoutRegex: `WebHook gave 500.*failed \d times to send the WebHook [^-]+-n$`, notifiers: shoutrrr.Slice{"fail": testNotifier(true, false)}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			stdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			webhook := testWebHook(false, true, false)
			if tc.wouldFail {
				webhook = testWebHook(true, true, false)
			}
			webhook.Delay = tc.delay
			maxTries := uint(tc.tries + 1)
			webhook.MaxTries = &maxTries
			webhook.SilentFails = &tc.silentFails
			webhook.Notifiers = &Notifiers{Shoutrrr: &tc.notifiers}
			if tc.tries > 0 {
				go func() {
					time.Sleep(time.Duration(11*(tc.tries-1)) * time.Second)
					webhook.Secret = "argus"
				}()
			}

			// WHEN try is called with it
			webhook.Send(utils.ServiceInfo{}, tc.useDelay)

			// THEN the logs are expected
			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = stdout
			output := string(out)
			re := regexp.MustCompile(tc.stdoutRegex)
			output = strings.ReplaceAll(output, "\n", "-n")
			match := re.MatchString(output)
			if !match {
				t.Errorf("%s:\nmatch on %q not found in\n%q",
					name, tc.stdoutRegex, output)
			}
		})
	}
}

func TestSliceSend(t *testing.T) {
	// GIVEN a Slice
	testLogging("INFO")
	tests := map[string]struct {
		slice          *Slice
		stdoutRegex    string
		stdoutRegexAlt string
		notifiers      shoutrrr.Slice
		useDelay       bool
		delays         map[string]string
		repeat         int
	}{
		"nil slice": {slice: nil, stdoutRegex: `^$`},
		"successful and failing webhook": {slice: &Slice{"pass": testWebHook(false, true, false), "fail": testWebHook(true, true, false)},
			stdoutRegex: `WebHook received.*failed \d times to send the WebHook`, stdoutRegexAlt: `failed \d times to send the WebHook.*WebHook received`},
		"does apply webhook delay": {slice: &Slice{"pass": testWebHook(false, true, false), "fail": testWebHook(true, true, false)},
			stdoutRegex: `WebHook received.*failed \d times to send the WebHook`, useDelay: true,
			delays: map[string]string{"fail": "2s", "pass": "1ms"}, repeat: 5},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.repeat++ // repeat to check delay usage as map order is random
			for tc.repeat != 0 {
				stdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w
				if tc.slice != nil {
					for id := range *tc.slice {
						(*tc.slice)[id].ID = id
					}
					for id := range tc.delays {
						(*tc.slice)[id].Delay = tc.delays[id]
					}
				}

				// WHEN try is called with it
				tc.slice.Send(utils.ServiceInfo{}, tc.useDelay)

				// THEN the logs are expected
				w.Close()
				out, _ := ioutil.ReadAll(r)
				os.Stdout = stdout
				output := string(out)
				output = strings.ReplaceAll(output, "\n", "-n")
				re := regexp.MustCompile(tc.stdoutRegex)
				match := re.MatchString(output)
				if !match {
					if tc.stdoutRegexAlt != "" {
						re = regexp.MustCompile(tc.stdoutRegexAlt)
						match = re.MatchString(output)
						if !match {
							t.Errorf("%s:\nmatch on %q not found in\n%q",
								name, tc.stdoutRegexAlt, output)
						}
						return
					}
					t.Errorf("%s:\nmatch on %q not found in\n%q",
						name, tc.stdoutRegex, output)
				}
				tc.repeat--
			}
		})
	}
}

func TestNotifiersSendWithNotifier(t *testing.T) {
	// GIVEN Notifiers
	testLogging("INFO")
	tests := map[string]struct {
		shoutrrrNotifiers *shoutrrr.Slice
		errRegex          string
	}{
		"nill Notifiers":      {errRegex: "^$"},
		"successful notifier": {errRegex: "^$", shoutrrrNotifiers: &shoutrrr.Slice{"pass": testNotifier(false, false)}},
		"failing notifier":    {errRegex: "invalid gotify token", shoutrrrNotifiers: &shoutrrr.Slice{"fail": testNotifier(true, false)}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			notifiers := Notifiers{Shoutrrr: tc.shoutrrrNotifiers}

			// WHEN Send is called with them
			err := notifiers.Send("TestNotifiersSendWithNotifier", name, &utils.ServiceInfo{})

			// THEN err is as expected
			e := utils.ErrorToString(err)
			re := regexp.MustCompile(tc.errRegex)
			match := re.MatchString(e)
			if !match {
				t.Errorf("%s:\nmatch on %q not found in\n%q",
					name, tc.errRegex, e)
			}
		})
	}
}

func TestCheckWebHookBody(t *testing.T) {
	// GIVEN a response body
	tests := map[string]struct {
		body string
		want bool
	}{
		"empty body":               {body: "", want: true},
		"success body":             {body: "success", want: true},
		"awx invalid secret":       {body: `{"detail":"You do not have permission to perform this action."}`, want: false},
		"adnanh/webhook hook fail": {body: `Hook rules were not satisfied.`, want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// WHEN checkWebHookBody is called on it
			got := checkWebHookBody(tc.body)

			// THEN the function returns the correct result
			if got != tc.want {
				t.Errorf("%s:\nwant: %t\ngot:  %t",
					name, tc.want, got)
			}
		})
	}
}
