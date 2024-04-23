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

package filter

import (
	"regexp"
	"strings"
	"testing"

	"github.com/release-argus/Argus/test"
	"github.com/release-argus/Argus/util"
	"gopkg.in/yaml.v3"
)

func TestURLCommandSlice_String(t *testing.T) {
	// GIVEN a URLCommandSlice
	tests := map[string]struct {
		slice *URLCommandSlice
		want  string
	}{
		"regex": {
			slice: &URLCommandSlice{
				testURLCommandRegex()},
			want: `
- type: regex
  regex: -([0-9.]+)-
`,
		},
		"regex (templated)": {
			slice: &URLCommandSlice{
				testURLCommandRegexTemplate()},
			want: `
- type: regex
  regex: -([0-9.]+)-
  template: _$1_
`,
		},
		"replace": {
			slice: &URLCommandSlice{
				testURLCommandReplace()},
			want: `
- type: replace
  new: bar
  old: foo
`,
		},
		"split": {
			slice: &URLCommandSlice{
				testURLCommandSplit()},
			want: `
- type: split
  index: 1
  text: this
`,
		},
		"all types": {
			slice: &URLCommandSlice{
				testURLCommandRegex(),
				testURLCommandReplace(),
				testURLCommandSplit()},
			want: `
- type: regex
  regex: -([0-9.]+)-
- type: replace
  new: bar
  old: foo
- type: split
  index: 1
  text: this
`,
		},
		"nil slice": {
			slice: nil,
			want:  "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			// WHEN String is called on it
			got := tc.slice.String()

			// THEN the expected string is returned
			tc.want = strings.TrimPrefix(tc.want, "\n")
			if got != tc.want {
				t.Fatalf("\nwant: %q\n got: %q",
					tc.want, got)
			}
		})
	}
}

func TestURLCommandsFromStr(t *testing.T) {
	// GIVEN a JSON string and a defaults URLCommandSlice
	defaults := URLCommandSlice{{Type: "regex"}}
	tests := map[string]struct {
		jsonStr  *string
		errRegex string
		want     *URLCommandSlice
	}{
		"regex - invalid": {
			jsonStr:  test.StringPtr(`[{"type":"regex","regex":"-([0-9.+)-"}]`),
			want:     &defaults,
			errRegex: `regex:.*\(Invalid RegEx\)`,
		},
		"regex": {
			jsonStr: test.StringPtr(`[{"type":"regex","regex":"-([0-9.]+)-"}]`),
			want: &URLCommandSlice{
				testURLCommandRegex()},
		},
		"replace": {
			jsonStr: test.StringPtr(`[{"type":"replace","old":"foo","new":"bar"}]`),
			want: &URLCommandSlice{
				testURLCommandReplace()},
		},
		"split": {
			jsonStr: test.StringPtr(`[{"type":"split","text":"this","index":1}]`),
			want: &URLCommandSlice{
				testURLCommandSplit()},
		},
		"all types": {
			jsonStr: test.StringPtr(`[{"type":"regex","regex":"-([0-9.]+)-"},{"type":"replace","old":"foo","new":"bar"},{"type":"split","text":"this","index":1}]`),
			want: &URLCommandSlice{
				testURLCommandRegex(),
				testURLCommandReplace(),
				testURLCommandSplit()},
		},
		"multiple of the each type": {
			jsonStr: test.StringPtr(`[{"type":"regex","regex":"-([0-9.]+)-"},{"type":"regex","regex":"-([0-9.]+)-"},{"type":"replace","old":"foo","new":"bar"},{"type":"replace","old":"foo","new":"bar"},{"type":"split","text":"this","index":1},{"type":"split","text":"this","index":1}]`),
			want: &URLCommandSlice{
				testURLCommandRegex(),
				testURLCommandRegex(),
				testURLCommandReplace(),
				testURLCommandReplace(),
				testURLCommandSplit(),
				testURLCommandSplit()},
		},
		"empty": {
			jsonStr: test.StringPtr(`[]`),
			want:    &URLCommandSlice{},
		},
		"object rather than list": {
			jsonStr:  test.StringPtr(`{"type":"regex"}`),
			errRegex: "cannot unmarshal object",
			want:     &defaults,
		},
		"nil": {
			jsonStr: nil,
			want:    &defaults,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN URLCommandsFromStr is called with it
			got, err := URLCommandsFromStr(tc.jsonStr, &defaults, &util.LogFrom{Primary: name})

			// THEN the expected URLCommandSlice is returned
			if err != nil {
				re := regexp.MustCompile(tc.errRegex)
				match := re.MatchString(err.Error())
				if !match {
					t.Errorf("want match for %q\nnot: %q",
						tc.errRegex, err)
				}
			}
			if got != tc.want {
				if got.String() != tc.want.String() {
					t.Fatalf("URLCommandsFromStr should have returned the expected URLCommandSlice:\nwant: %q\ngot:  %q",
						tc.want, got)
				}
			}
		})
	}
}

func TestLogInit(t *testing.T) {
	// GIVEN a JLog
	newJLog := util.NewJLog("WARN", false)

	// WHEN LogInit is called with it
	LogInit(newJLog)

	// THEN the global JLog is set to its address
	if jLog != newJLog {
		t.Fatalf("JLog should have been initialised to the one we called Init with")
	}
}

func TestURLCommandSlice_Run(t *testing.T) {
	// GIVEN a URLCommandSlice
	testText := "abc123-def456"
	tests := map[string]struct {
		slice    *URLCommandSlice
		text     string
		want     string
		errRegex string
	}{
		"nil slice": {
			slice:    nil,
			errRegex: "^$",
			want:     testText,
		},
		"regex": {
			slice: &URLCommandSlice{
				{Type: "regex", Regex: test.StringPtr("([a-z]+)[0-9]+"), Index: 1}},
			errRegex: "^$",
			want:     "def",
		},
		"regex with negative index": {
			slice: &URLCommandSlice{
				{Type: "regex", Regex: test.StringPtr("([a-z]+)[0-9]+"), Index: -1}},
			errRegex: "^$",
			want:     "def",
		},
		"regex doesn't match (gives text that didn't match)": {
			slice: &URLCommandSlice{
				{Type: "regex", Regex: test.StringPtr("([h-z]+)[0-9]+"), Index: 1}},
			errRegex: `regex .* didn't return any matches on "` + testText + `"`,
			want:     testText,
		},
		"regex doesn't match (doesn't give text that didn't match as too long)": {
			slice: &URLCommandSlice{
				{Type: "regex", Regex: test.StringPtr("([h-z]+)[0-9]+"), Index: 1}},
			errRegex: "regex .* didn't return any matches$",
			text:     strings.Repeat("a123", 5),
			want:     "a123a123a123a123a123",
		},
		"regex index out of bounds": {
			slice: &URLCommandSlice{
				{Type: "regex", Regex: test.StringPtr("([a-z]+)[0-9]+"), Index: 2}},
			errRegex: `regex .* returned \d elements on "[^']+", but the index wants element number \d`,
			want:     testText,
		},
		"regex with template": {
			slice: &URLCommandSlice{
				{Type: "regex", Regex: test.StringPtr("([a-z]+)([0-9]+)"), Index: 1, Template: test.StringPtr("$1_$2")}},
			errRegex: "^$",
			want:     "def_456",
		},
		"replace": {
			slice: &URLCommandSlice{
				{Type: "replace", Old: test.StringPtr("-"), New: test.StringPtr(" ")}},
			errRegex: "^$",
			want:     "abc123 def456",
		},
		"split": {
			slice: &URLCommandSlice{
				{Type: "split", Text: test.StringPtr("-"), Index: -1}},
			errRegex: "^$",
			want:     "def456",
		},
		"split with negative index": {
			slice: &URLCommandSlice{
				{Type: "split", Text: test.StringPtr("-"), Index: 0}},
			errRegex: "^$",
			want:     "abc123",
		},
		"split on unknown text": {
			slice: &URLCommandSlice{
				{Type: "split", Text: test.StringPtr("7"), Index: 0}},
			errRegex: "split didn't find any .* to split on",
			want:     testText,
		},
		"split index out of bounds": {
			slice: &URLCommandSlice{
				{Type: "split", Text: test.StringPtr("-"), Index: 2}},
			errRegex: `split .* returned \d elements on "[^']+", but the index wants element number \d`,
			want:     testText,
		},
		"all types": {
			slice: &URLCommandSlice{
				{Type: "regex", Regex: test.StringPtr("([a-z]+)[0-9]+"), Index: 1},
				{Type: "replace", Old: test.StringPtr("e"), New: test.StringPtr("a")},
				{Type: "split", Text: test.StringPtr("a"), Index: 1}},
			errRegex: "^$",
			want:     "f",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN run is called on it
			text := testText
			if tc.text != "" {
				text = tc.text
			}
			text, err := tc.slice.Run(text, &util.LogFrom{})

			// THEN the expected text was returned
			if tc.want != text {
				t.Errorf("Should have got %q, not %q",
					tc.want, text)
			}
			e := util.ErrorToString(err)
			re := regexp.MustCompile(tc.errRegex)
			match := re.MatchString(e)
			if !match {
				t.Fatalf("want match for %q\nnot: %q",
					tc.errRegex, e)
			}
		})
	}
}

func TestURLCommand_String(t *testing.T) {
	// GIVEN a URLCommand
	regex := testURLCommandRegex()
	replace := testURLCommandReplace()
	split := testURLCommandSplit()
	tests := map[string]struct {
		cmd  *URLCommand
		want string
	}{
		"regex": {
			cmd: &regex,
			want: `
type: regex
regex: -([0-9.]+)-
`,
		},
		"replace": {
			cmd: &replace,
			want: `
type: replace
new: bar
old: foo
`,
		},
		"split": {
			cmd: &split,
			want: `
type: split
index: 1
text: this
`,
		},
		"nil slice": {
			cmd:  nil,
			want: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN String is called on it
			got := tc.cmd.String()

			// THEN the expected string is returned
			tc.want = strings.TrimPrefix(tc.want, "\n")
			if got != tc.want {
				t.Fatalf("\nwant: %q\n got: %q",
					tc.want, got)
			}
		})
	}
}

func TestURLCommandSlice_CheckValues(t *testing.T) {
	// GIVEN a URLCommandSlice
	tests := map[string]struct {
		slice     *URLCommandSlice
		wantSlice *URLCommandSlice
		errRegex  []string
	}{
		"nil slice": {
			slice:    nil,
			errRegex: []string{`^$`},
		},
		"valid regex": {
			slice:    &URLCommandSlice{testURLCommandRegex()},
			errRegex: []string{`^$`},
		},
		"undefined regex": {
			slice: &URLCommandSlice{
				{Type: "regex"}},
			errRegex: []string{`^url_commands:$`, `^  item_0:$`, `^    regex: <required>`},
		},
		"invalid regex": {
			slice: &URLCommandSlice{
				{Type: "regex", Regex: test.StringPtr("[0-")}},
			errRegex: []string{`^    regex: .* <invalid>`},
		},
		"valid regex with template": {
			slice:    &URLCommandSlice{testURLCommandRegexTemplate()},
			errRegex: []string{`^$`},
		},
		"valid regex with empty template": {
			slice: &URLCommandSlice{
				{Type: "regex", Regex: test.StringPtr("[0-"), Template: test.StringPtr("")}},
			wantSlice: &URLCommandSlice{
				{Type: "regex", Regex: test.StringPtr("[0-")}},
			errRegex: []string{`^$`},
		},
		"valid replace": {
			slice: &URLCommandSlice{
				testURLCommandReplace()},
			errRegex: []string{`^$`},
		},
		"invalid replace": {
			slice: &URLCommandSlice{
				{Type: "replace"}},
			errRegex: []string{`^    new: <required>`, `^    old: <required>`},
		},
		"valid split": {
			slice: &URLCommandSlice{
				testURLCommandSplit()},
			errRegex: []string{`^$`},
		},
		"invalid split": {
			slice: &URLCommandSlice{
				{Type: "split"}},
			errRegex: []string{`^    text: <required>`},
		},
		"invalid type": {
			slice: &URLCommandSlice{
				{Type: "something"}},
			errRegex: []string{`^    type: .* <invalid>`},
		},
		"valid all types": {
			slice: &URLCommandSlice{
				testURLCommandRegex(),
				testURLCommandReplace(),
				testURLCommandSplit()},
			errRegex: []string{`^$`},
		},
		"all possible errors": {
			slice: &URLCommandSlice{
				{Type: "regex"}, {Type: "replace"},
				{Type: "split"},
				{Type: "something"}},
			errRegex: []string{
				`^url_commands:$`,
				`^  item_0:$`,
				`^    regex: <required>`,
				`^    new: <required>`,
				`^    old: <required>`,
				`^    text: <required>`,
				`^    type: .* <invalid>`},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN CheckValues is called on it
			err := tc.slice.CheckValues("")

			// THEN err is expected
			e := util.ErrorToString(err)
			lines := strings.Split(e, `\`)
			for i := range tc.errRegex {
				re := regexp.MustCompile(tc.errRegex[i])
				found := false
				for j := range lines {
					match := re.MatchString(lines[j])
					if match {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("want match for: %q\ngot:  %q",
						tc.errRegex[i], strings.ReplaceAll(e, `\`, "\n"))
				}
			}

			// AND the slice is as expected
			if tc.wantSlice != nil {
				strHave := tc.slice.String()
				strWant := tc.wantSlice.String()
				if strHave != strWant {
					t.Errorf("want slice:\n%q\ngot:  %q",
						strWant, strHave)
				}
			}
		})
	}
}

func TestURLCommandSlice_UnmarshalYAML(t *testing.T) {
	// GIVEN a file to read a URLCommandSlice
	tests := map[string]struct {
		input    string
		slice    URLCommandSlice
		errRegex string
	}{
		"invalid unmarshal": {
			input: `type: regex
regex: foo
regex: foo
index: 1
text: hi
old: was
new: now`,
			errRegex: "mapping key .* already defined",
		},
		"non-list URLCommand": {
			input: `type: regex
regex: foo
index: 1
text: hi
old: was
new: now`,
			slice: URLCommandSlice{
				{Type: "regex",
					Regex: test.StringPtr("foo"), Index: 1,
					Text: test.StringPtr("hi"), Old: test.StringPtr("was"), New: test.StringPtr("now")}},
			errRegex: "^$",
		},
		"list of URLCommands": {
			input: `- type: regex
  regex: \"([0-9.+])\"
  index: 1
- type: replace
  old: foo
  new: bar
- type: split
  text: abc
  index: 2`,
			errRegex: "^$",
			slice: URLCommandSlice{
				{Type: "regex",
					Regex: test.StringPtr(`\"([0-9.+])\"`), Index: 1},
				{Type: "replace", Old: test.StringPtr("foo"), New: test.StringPtr("bar")},
				{Type: "split", Text: test.StringPtr("abc"), Index: 2}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			var slice URLCommandSlice

			// WHEN Unmarshalled
			err := yaml.Unmarshal([]byte(tc.input), &slice)

			// THEN the it errs when appropriate and unmarshals correctly into a list
			e := util.ErrorToString(err)
			re := regexp.MustCompile(tc.errRegex)
			match := re.MatchString(e)
			if !match {
				t.Fatalf("want match for %q\nnot: %q",
					tc.errRegex, e)
			}
			if len(slice) != len(tc.slice) {
				t.Fatalf("got a slice of length %d. want %d\n%#v",
					len(slice), len(tc.slice), slice)
			}
			for i := range tc.slice {
				if slice[i].Type != tc.slice[i].Type {
					t.Errorf("wrong Type:\nwant: %q\ngot:  %q\n",
						tc.slice[i].Type, slice[i].Type)
				}
				if util.DefaultIfNil(slice[i].Regex) != util.DefaultIfNil(tc.slice[i].Regex) {
					t.Errorf("wrong Regex:\nwant: %q\ngot:  %q\n",
						util.DefaultIfNil(tc.slice[i].Regex), util.DefaultIfNil(slice[i].Regex))
				}
				if slice[i].Index != tc.slice[i].Index {
					t.Errorf("wrong Index:\nwant: %q\ngot:  %q\n",
						tc.slice[i].Index, slice[i].Index)
				}
				if util.DefaultIfNil(slice[i].Text) != util.DefaultIfNil(tc.slice[i].Text) {
					t.Errorf("wrong Text:\nwant: %q\ngot:  %q\n",
						util.DefaultIfNil(tc.slice[i].Text), util.DefaultIfNil(slice[i].Text))
				}
				if util.DefaultIfNil(slice[i].Old) != util.DefaultIfNil(tc.slice[i].Old) {
					t.Errorf("wrong Old:\nwant: %q\ngot:  %q\n",
						util.DefaultIfNil(tc.slice[i].Old), util.DefaultIfNil(slice[i].Old))
				}
				if util.DefaultIfNil(slice[i].New) != util.DefaultIfNil(tc.slice[i].New) {
					t.Errorf("wrong New:\nwant: %q\ngot:  %q\n",
						util.DefaultIfNil(tc.slice[i].New), util.DefaultIfNil(slice[i].New))
				}
			}
		})
	}
}
