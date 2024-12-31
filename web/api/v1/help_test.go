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

package v1

import (
	"encoding/json"
	"os"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/release-argus/Argus/command"
	"github.com/release-argus/Argus/config"
	dbtype "github.com/release-argus/Argus/db/types"
	"github.com/release-argus/Argus/notify/shoutrrr"
	"github.com/release-argus/Argus/service"
	deployedver "github.com/release-argus/Argus/service/deployed_version"
	latestver "github.com/release-argus/Argus/service/latest_version"
	opt "github.com/release-argus/Argus/service/option"
	"github.com/release-argus/Argus/test"
	"github.com/release-argus/Argus/util"
	"github.com/release-argus/Argus/webhook"
)

var (
	loadMutex             sync.Mutex
	loadCount             int
	secretValueMarshalled string
)

func TestMain(m *testing.M) {
	// initialise jLog
	masterJLog := util.NewJLog("DEBUG", false)
	masterJLog.Testing = true
	flags := make(map[string]bool)
	path := "TestWebAPIv1Main.yml"
	testYAML_Argus(path)
	var config config.Config
	config.Load(path, &flags, masterJLog)
	os.Remove(path)
	LogInit(masterJLog)

	// Marshal the secret value '<secret>' -> '\u003csecret\u003e'
	secretValueMarshalledBytes, _ := json.Marshal(util.SecretValue)
	secretValueMarshalled = string(secretValueMarshalledBytes)

	// run other tests
	exitCode := m.Run()

	// exit
	os.Exit(exitCode)
}

func testClient() Client {
	hub := NewHub()
	return Client{
		hub:  hub,
		ip:   "1.1.1.1",
		conn: &websocket.Conn{},
		send: make(chan []byte, 5),
	}
}

func testLoad(file string) *config.Config {
	var config config.Config

	flags := make(map[string]bool)
	config.Load(file, &flags, nil)
	config.Init(false)
	announceChannel := make(chan []byte, 8)
	config.HardDefaults.Service.Status.AnnounceChannel = &announceChannel

	return &config
}

func testAPI(name string) API {
	testYAML_Argus(name)

	cfg := testLoad(name)
	cfg.HardDefaults.Service.LatestVersion.AccessToken = os.Getenv("GITHUB_TOKEN")
	return API{Config: cfg}
}

func testService(id string) *service.Service {
	announceChannel := make(chan []byte, 8)
	databaseChannel := make(chan dbtype.Message, 8)

	lv, _ := latestver.New(
		"url",
		"yaml", test.TrimYAML(`
			url: https://valid.release-argus.io/plain
			url_commands:
				- type: regex
					regex: 'stable version: "v?([0-9.]+)"'
		`),
		nil,
		nil,
		nil, nil)

	dv, _ := deployedver.New(
		"yaml", test.TrimYAML(`
			method: GET
			url: https://valid.release-argus.io/json
			json: foo.bar.version
		`),
		nil,
		nil,
		nil, nil)

	options := opt.New(
		nil, "", test.BoolPtr(true),
		nil, nil)

	svc := service.Service{
		ID:                    id,
		Comment:               "foo",
		LatestVersion:         lv,
		DeployedVersionLookup: dv,
		Options:               *options}

	// Hard defaults
	serviceHardDefaults := service.Defaults{}
	serviceHardDefaults.Default()
	shoutrrrHardDefaults := shoutrrr.SliceDefaults{}
	shoutrrrHardDefaults.Default()
	webhookHardDefaults := webhook.Defaults{}
	webhookHardDefaults.Default()

	// Defaults
	serviceDefaults := service.Defaults{}
	serviceDefaults.Init()

	// Init with defaults/hardDefaults
	svc.Init(
		&serviceDefaults, &serviceHardDefaults,
		&shoutrrr.SliceDefaults{}, &shoutrrr.SliceDefaults{}, &shoutrrrHardDefaults,
		&webhook.SliceDefaults{}, &webhook.Defaults{}, &webhookHardDefaults)

	// Status channels
	svc.Status.AnnounceChannel = &announceChannel
	svc.Status.DatabaseChannel = &databaseChannel

	return &svc
}

func testCommand(failing bool) command.Command {
	if failing {
		return command.Command{"ls", "-lah", "/root"}
	}
	return command.Command{"ls", "-lah"}
}

func testFaviconSettings(png string, svg string) *config.FaviconSettings {
	if svg == "" && png == "" {
		return nil
	}

	return &config.FaviconSettings{
		SVG: svg,
		PNG: png}
}
