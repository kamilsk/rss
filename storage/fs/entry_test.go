package fs_test

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/kamilsk/stream/storage/fs"
)

var Update = flag.Bool("update", false, "update .golden files")

func TestEntry_MarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		entry    fs.Entry
		expected struct {
			err    error
			golden string
		}
	}{
		{
			entry: fs.Entry{ID: "end_resource", Name: "end resource"},
			expected: struct {
				err    error
				golden string
			}{golden: "./fixtures/end_resource.golden"},
		},
	} {
		data, err := json.Marshal(&tc.entry)
		switch {
		case tc.expected.err != nil && err != nil:
			if tc.expected.err.Error() != err.Error() {
				t.Errorf("unexpected error. expected: %v, obtained: %v", tc.expected.err, err)
			}
		case tc.expected.err == nil && err != nil:
			fallthrough
		case tc.expected.err != nil && err == nil:
			t.Errorf("unexpected error. expected: %v, obtained: %v", tc.expected.err, err)
		}
		if true || *Update {
			json.NewEncoder(func(file string) io.Writer {
				f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				if err != nil {
					panic(err)
				}
				return f
			}(tc.expected.golden)).Encode(tc.entry)
		}
		golden := func(file string) []byte {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}
			return data
		}(tc.expected.golden)
		obtained, expected := strings.TrimSpace(string(data)), strings.TrimSpace(string(golden))
		if expected != obtained {
			t.Errorf("unexpected JSON. expected: %q, obtained: %q", expected, obtained)
		}
	}
}

func TestEntry_UnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		file     string
		expected struct {
			err     error
			visitor func(fs.Entry) bool
		}
	}{
		{
			file: "./fixtures/end_resource.json",
			expected: struct {
				err     error
				visitor func(fs.Entry) bool
			}{
				visitor: func(e fs.Entry) bool {
					return e.URL() == "https://rss.octolab.net/kamilsk/podcasts"
				},
			},
		},
		{
			file: "./fixtures/multi_resource.json",
			expected: struct {
				err     error
				visitor func(fs.Entry) bool
			}{
				visitor: func(e fs.Entry) bool {
					return len(e.SubEntries()) == 3 && func() bool {
						for i, src := range map[int]string{
							0: "https://rss.octolab.net/kamilsk/podcasts",
							1: "https://rss.octolab.net/kamilsk/releases",
							2: "https://rss.octolab.net/octolab/releases",
						} {
							if e.SubEntries()[i].URL() != src {
								return false
							}
						}
						return true
					}()
				},
			},
		},
		{
			file: "./fixtures/mixed_resource.json",
			expected: struct {
				err     error
				visitor func(fs.Entry) bool
			}{
				visitor: func(e fs.Entry) bool {
					return len(e.SubEntries()) == 2 &&
						e.SubEntries()[1].URL() == "https://rss.octolab.net/octolab/releases" && func() bool {
						for i, src := range map[int]string{
							0: "https://rss.octolab.net/kamilsk/podcasts",
							1: "https://rss.octolab.net/kamilsk/releases",
						} {
							if e.SubEntries()[0].SubEntries()[i].URL() != src {
								return false
							}
						}
						return true
					}()
				},
			},
		},
	} {
		var obtained fs.Entry
		err := json.NewDecoder(func(file string) io.Reader {
			f, err := os.Open(file)
			if err != nil {
				panic(err)
			}
			return f
		}(tc.file)).Decode(&obtained)
		switch {
		case tc.expected.err != nil && err != nil:
			if tc.expected.err.Error() != err.Error() {
				t.Errorf("unexpected error. expected: %v, obtained: %v", tc.expected.err, err)
			}
		case tc.expected.err == nil && err != nil:
			fallthrough
		case tc.expected.err != nil && err == nil:
			t.Errorf("unexpected error. expected: %v, obtained: %v", tc.expected.err, err)
		}
		if !tc.expected.visitor(obtained) {
			t.Errorf("visitor failed at ID:%s %q", obtained.ID, obtained.Name)
		}
	}
}
