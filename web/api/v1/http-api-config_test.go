// Copyright [2023] [Argus]
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

package v1

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	command "github.com/release-argus/Argus/commands"
	"github.com/release-argus/Argus/config"
	"github.com/release-argus/Argus/notifiers/shoutrrr"
	"github.com/release-argus/Argus/service"
	latestver "github.com/release-argus/Argus/service/latest_version"
	"github.com/release-argus/Argus/service/latest_version/filter"
	opt "github.com/release-argus/Argus/service/options"
	"github.com/release-argus/Argus/test"
	"github.com/release-argus/Argus/webhook"
)

func TestHTTP_Config(t *testing.T) {
	// GIVEN an API and a request for the config
	file := "TestHTTP_Config.yml"
	api := testAPI(file)
	defer func() {
		os.RemoveAll(file)
		if api.Config.Settings.Data.DatabaseFile != nil {
			os.RemoveAll(*api.Config.Settings.Data.DatabaseFile)
		}
	}()

	tests := map[string]struct {
		settings *config.Settings
		defaults *config.Defaults
		notify   *shoutrrr.SliceDefaults
		webhook  *webhook.SliceDefaults
		service  *service.Slice
		order    *[]string
		wantBody string
	}{
		"0. settings": {
			settings: &config.Settings{
				SettingsBase: config.SettingsBase{
					Web: config.WebSettings{
						ListenHost: test.StringPtr("127.0.0.1")}}},
			wantBody: `
				{
					"settings": {
						"log": {},
						"web": {
							"listen_host": "127.0.0.1"
						}
					},
					"defaults": {
						"service": {
							"options": {},
							"latest_version": {
								"require": {
									"docker": {}
								}
							},
							"deployed_version": {},
							"dashboard": {}
						},
						"webhook": {}
					},
					"notify": {},
					"webhook": {},
					"service": {}
				}`,
		},
		"1. settings + defaults": {
			defaults: &config.Defaults{
				Service: service.Defaults{
					Options: opt.OptionsDefaults{
						OptionsBase: opt.OptionsBase{
							Interval: "1h"}},
					LatestVersion: *latestver.NewDefaults(
						test.StringPtr("foo"),
						test.BoolPtr(true),
						test.BoolPtr(false),
						filter.NewRequireDefaults(
							filter.NewDockerCheckDefaults(
								"ghcr",
								"tokenForGHCR",
								"tokenForHub", "usernameForHub",
								"tokenForQuay",
								nil)))}},
			wantBody: `
				{
					"settings": {
						"log": {},
						"web": {
							"listen_host": "127.0.0.1"
						}
					},
					"defaults": {
						"service": {
							"options": {
								"interval": "1h"
							},
							"latest_version": {
								"access_token": "\u003csecret\u003e",
								"allow_invalid_certs": true,
								"use_prerelease": false,
								"require": {
									"docker": {
										"type": "ghcr",
										"ghcr": {
											"token": "\u003csecret\u003e"
										},
										"hub": {
											"token": "\u003csecret\u003e",
											"username": "usernameForHub"
										},
										"quay": {
											"token": "\u003csecret\u003e"
										}
									}
								}
							},
							"deployed_version": {},
							"dashboard": {}
						},
						"webhook": {}
					},
					"notify": {},
					"webhook": {},
					"service": {}
				}`,
		},
		"2. settings + defaults (with notify+command+webhook service defaults)": {
			defaults: &config.Defaults{
				Service: service.Defaults{
					Options: opt.OptionsDefaults{
						OptionsBase: opt.OptionsBase{
							Interval: "1h"}},
					LatestVersion: *latestver.NewDefaults(
						test.StringPtr("foo"),
						test.BoolPtr(true),
						test.BoolPtr(false),
						filter.NewRequireDefaults(
							filter.NewDockerCheckDefaults(
								"ghcr",
								"tokenForGHCR",
								"tokenForHub", "usernameForHub",
								"tokenForQuay",
								nil))),
					Notify: map[string]struct{}{
						"n1": {}},
					Command: command.Slice{
						{"command", "arg1", "arg2"}},
					WebHook: map[string]struct{}{
						"wh1": {},
						"wh2": {},
						"wh3": {},
						"wh4": {}}}},
			wantBody: `
				{
					"settings": {
						"log": {},
						"web": {
							"listen_host": "127.0.0.1"
						}
					},
					"defaults": {
						"service": {
							"options": {
								"interval": "1h"
							},
							"latest_version": {
								"access_token": "\u003csecret\u003e",
								"allow_invalid_certs": true,
								"use_prerelease": false,
								"require": {
									"docker": {
										"type": "ghcr",
										"ghcr": {
											"token": "\u003csecret\u003e"
										},
										"hub": {
											"token": "\u003csecret\u003e",
											"username": "usernameForHub"
										},
										"quay": {
											"token": "\u003csecret\u003e"
										}
									}
								}
							},
							"notify": {
								"n1": {}
							},
							"command": [
								["command","arg1","arg2"]
							],
							"webhook": {
								"wh1": {},
								"wh2": {},
								"wh3": {},
								"wh4": {}
							},
							"deployed_version": {},
							"dashboard": {}
						},
						"webhook": {}
					},
					"notify": {},
					"webhook": {},
					"service": {}
				}`,
		},
		"3. settings + defaults (with notify+command+webhook service defaults) + notify": {
			notify: &shoutrrr.SliceDefaults{
				"foo": shoutrrr.NewDefaults(
					"gotify",
					&map[string]string{
						"message": "hello world"},
					&map[string]string{
						"title": "UPDATE"},
					&map[string]string{
						"token": "secret123"})},
			wantBody: `
				{
					"settings": {
						"log": {},
						"web": {
							"listen_host": "127.0.0.1"
						}
					},
					"defaults": {
						"service": {
							"options": {
								"interval": "1h"
							},
							"latest_version": {
								"access_token": "\u003csecret\u003e",
								"allow_invalid_certs": true,
								"use_prerelease": false,
								"require": {
									"docker": {
										"type": "ghcr",
										"ghcr": {
											"token": "\u003csecret\u003e"
										},
										"hub": {
											"token": "\u003csecret\u003e",
											"username": "usernameForHub"
										},
										"quay": {
											"token": "\u003csecret\u003e"
										}
									}
								}
							},
							"notify": {
								"n1": {}
							},
							"command": [
								["command","arg1","arg2"]
							],
							"webhook": {
								"wh1": {},
								"wh2": {},
								"wh3": {},
								"wh4": {}
							},
							"deployed_version": {},
							"dashboard": {}
						},
						"webhook": {}
					},
					"notify": {
						"foo": {
							"type": "gotify",
							"options": {
								"message": "hello world"
							},
							"url_fields": {
								"token": "\u003csecret\u003e"
							},
							"params": {
								"title": "UPDATE"
							}
						}
					},
					"webhook": {},
					"service": {}
				}`,
		},
		"4. settings + defaults (with notify+command+webhook service defaults) + notify + webhook": {
			webhook: &webhook.SliceDefaults{
				"foo": webhook.NewDefaults(
					test.BoolPtr(true), // allow_invalid_certs
					&webhook.Headers{
						{Key: "X-Header", Value: "value"}},
					"4s",                    // delay
					nil,                     // desired_status_code
					nil,                     // max_tries
					"something",             // secret
					nil,                     // silent_fails
					"github",                // type
					"https://example.com")}, // url
			wantBody: `
				{
					"settings": {
						"log": {},
						"web": {
							"listen_host": "127.0.0.1"
						}
					},
					"defaults": {
						"service": {
							"options": {
								"interval": "1h"
							},
							"latest_version": {
								"access_token": "\u003csecret\u003e",
								"allow_invalid_certs": true,
								"use_prerelease": false,
								"require": {
									"docker": {
										"type": "ghcr",
										"ghcr": {
											"token": "\u003csecret\u003e"
										},
										"hub": {
											"token": "\u003csecret\u003e",
											"username": "usernameForHub"
										},
										"quay": {
											"token": "\u003csecret\u003e"
										}
									}
								}
							},
							"notify": {
								"n1": {}
							},
							"command": [
								["command","arg1","arg2"]
							],
							"webhook": {
								"wh1": {},
								"wh2": {},
								"wh3": {},
								"wh4": {}
							},
							"deployed_version": {},
							"dashboard": {}
						},
						"webhook": {}
					},
					"notify": {
						"foo": {
							"type": "gotify",
							"options": {
								"message": "hello world"
							},
							"url_fields": {
								"token": "\u003csecret\u003e"
							},
							"params": {
								"title": "UPDATE"
							}
						}
					},
					"webhook": {
						"foo": {
							"type": "github",
							"url": "https://example.com",
							"allow_invalid_certs": true,
							"secret": "\u003csecret\u003e",
							"custom_headers": [
								{"key": "X-Header","value": "\u003csecret\u003e"}],
							"delay": "4s"}},
					"service": {}
				}`,
		},
		"5. settings + defaults (with notify+command+webhook service defaults) + notify + webhook + service": {
			service: &service.Slice{
				"alpha": &service.Service{
					LatestVersion: *latestver.New(
						test.StringPtr("aToken"),
						nil, nil, nil, nil, nil,
						"github",
						"release-argus/Argus",
						nil, nil, nil, nil)},
				"bravo": &service.Service{
					LatestVersion: *latestver.New(
						test.StringPtr("otherToken"),
						nil, nil, nil, nil, nil,
						"url",
						"https://example.com/version",
						nil, nil, nil, nil)}},
			order: &[]string{"alpha", "bravo"},
			wantBody: `
				{
					"settings": {
						"log": {},
						"web": {
							"listen_host": "127.0.0.1"
						}
					},
					"defaults": {
						"service": {
							"options": {
								"interval": "1h"
							},
							"latest_version": {
								"access_token": "\u003csecret\u003e",
								"allow_invalid_certs": true,
								"use_prerelease": false,
								"require": {
									"docker": {
										"type": "ghcr",
										"ghcr": {
											"token": "\u003csecret\u003e"
										},
										"hub": {
											"token": "\u003csecret\u003e",
											"username": "usernameForHub"
										},
										"quay": {
											"token": "\u003csecret\u003e"
										}
									}
								}
							},
							"notify": {
								"n1": {}
							},
							"command": [
								["command","arg1","arg2"]
							],
							"webhook": {
								"wh1": {},
								"wh2": {},
								"wh3": {},
								"wh4": {}
							},
							"deployed_version": {},
							"dashboard": {}
						},
						"webhook": {}
					},
					"notify": {
						"foo": {
							"type": "gotify",
							"options": {
								"message": "hello world"
							},
							"url_fields": {
								"token": "\u003csecret\u003e"
							},
							"params": {
								"title": "UPDATE"
							}
						}
					},
					"webhook": {
						"foo": {
							"type": "github",
							"url": "https://example.com",
							"allow_invalid_certs": true,
							"secret": "\u003csecret\u003e",
							"custom_headers": [
								{"key": "X-Header","value": "\u003csecret\u003e"}],
							"delay": "4s"}},
					"service": {
						"alpha": {
							"options": {},
							"latest_version": {
								"type": "github",
								"url": "release-argus/Argus",
								"access_token": "\u003csecret\u003e",
								"url_commands": []
							},
							"command": [],
							"notify": {},
							"webhook": {},
							"dashboard": {}
						},
						"bravo": {
							"options": {},
							"latest_version": {
								"type": "url",
								"url": "https://example.com/version",
								"access_token": "\u003csecret\u003e",
								"url_commands": []
							},"command": [],
							"notify": {},
							"webhook": {},
							"dashboard": {}
						}
					},
					"order": [
						"alpha",
						"bravo"
					]
				}`,
		},
	}
	order := make([]string, len(tests))
	for i := range order {
		lookingFor := fmt.Sprintf("%d. ", i)
		for name := range tests {
			if strings.HasPrefix(name, lookingFor) {
				order[i] = name
				break
			}
		}
	}
	api.Config.Settings = config.Settings{}
	api.Config.Defaults = config.Defaults{}
	api.Config.Notify = shoutrrr.SliceDefaults{}
	api.Config.WebHook = webhook.SliceDefaults{}
	api.Config.Service = service.Slice{}
	api.Config.Order = []string{}

	for _, name := range order {
		tc := tests[name]
		t.Run(name, func(t *testing.T) {

			if tc.settings != nil {
				api.Config.Settings = *tc.settings
			}
			if tc.defaults != nil {
				api.Config.Defaults = *tc.defaults
			}
			if tc.notify != nil {
				api.Config.Notify = *tc.notify
			}
			if tc.webhook != nil {
				api.Config.WebHook = *tc.webhook
			}
			if tc.service != nil {
				api.Config.Service = *tc.service
			}
			if tc.order != nil {
				api.Config.Order = *tc.order
			}
			tc.wantBody = test.TrimJSON(tc.wantBody) + "\n"

			// WHEN that HTTP request is sent
			req := httptest.NewRequest(http.MethodGet, "/api/v1/config", nil)
			w := httptest.NewRecorder()
			api.httpConfig(w, req)
			res := w.Result()
			defer res.Body.Close()

			// THEN the expected body is returned as expected
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("unexpected error - %v",
					err)
			}
			got := string(data)
			if got != tc.wantBody {
				t.Fatalf("want %q\ngot: %q",
					tc.wantBody, got)
			}
		})
	}
}
