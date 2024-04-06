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
	"regexp"
	"strings"
	"sync"
	"testing"
)

func TestHTTP_httpFlags(t *testing.T) {
	// GIVEN an API and a request for the flag var values
	file := "TestHTTP_httpFlags.yml"
	api := testAPI(file)
	apiMutex := sync.RWMutex{}
	defer func() {
		os.RemoveAll(file)
		if api.Config.Settings.Data.DatabaseFile != nil {
			os.RemoveAll(*api.Config.Settings.Data.DatabaseFile)
		}
	}()
	want := `
		{
			"config.file":"` + file + `",
			"log.level":"` + fmt.Sprintf(api.Config.Settings.LogLevel()) + `",
			"log.timestamps":` + fmt.Sprint(*api.Config.Settings.LogTimestamps()) + `,
			"data.database-file":"` + api.Config.Settings.DataDatabaseFile() + `",
			"web.listen-host":"` + api.Config.Settings.WebListenHost() + `",
			"web.listen-port":"[0-9]{1,5}",
			"web.cert-file":null,
			"web.pkey-file":null,
			"web.route-prefix":"` + strings.ReplaceAll(api.Config.Settings.WebRoutePrefix(), "/", `\/`) + `"
		}\s$`

	// WHEN that HTTP request is sent
	req := httptest.NewRequest(http.MethodGet, "/api/v1/flags", nil)
	w := httptest.NewRecorder()
	apiMutex.RLock()
	api.httpFlags(w, req)
	apiMutex.RUnlock()
	res := w.Result()
	defer res.Body.Close()

	// THEN the expected body is returned as expected
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("unexpected error - %v",
			err)
	}
	got := string(data)
	want = strings.ReplaceAll(want, "\t", "")
	want = strings.ReplaceAll(want, "\n", "")
	wantRe := regexp.MustCompile(want)
	if !wantRe.MatchString(got) {
		t.Errorf("%q doesn't match regex %q",
			got, want)
	}
}
