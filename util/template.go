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

package util

import (
	"strings"
	"sync"

	"github.com/flosch/pongo2/v5"
)

//	github.com/flosch/pongo2/v5.(*TemplateSet).FromString()
//
// /Users/josephkav/Git/GitHub/Argus/vendor/github.com/flosch/pongo2/v5/template_sets.go:199 +0x4f
// github.com/flosch/pongo2/v5.(*TemplateSet).FromString-fm()
//
//	<autogenerated>:1 +0x4d
var pongoMutex = sync.Mutex{}

// TemplateString with pongo2 and `context`.
func TemplateString(template string, context ServiceInfo) (result string) {
	// If the string isn't a Jinja template
	if !strings.Contains(template, "{") {
		result = template
		return
	}
	// pongo2 DATA RACE
	pongoMutex.Lock()
	defer pongoMutex.Unlock()

	// Compile the template.
	tpl, err := pongo2.FromString(template)
	if err != nil {
		panic(err)
	}

	// Render the template.
	result, err = tpl.Execute(pongo2.Context{
		"service_id":  context.ID,
		"service_url": context.URL,
		"web_url":     context.WebURL,
		"version":     context.LatestVersion,
	})
	if err != nil {
		panic(err)
	}
	return
}

// CheckTemplate will compile
//
// true == pass
//
// false == fail
func CheckTemplate(template string) (pass bool) {
	// pongo2 DATA RACE
	pongoMutex.Lock()
	defer pongoMutex.Unlock()

	_, err := pongo2.FromString(template)
	return err == nil
}
