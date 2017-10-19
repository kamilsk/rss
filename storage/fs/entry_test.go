package fs_test

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/kamilsk/stream/storage/fs"
)

func TestEntry_MarshalJSON(t *testing.T) {
	//
}

func TestEntry_UnmarshalJSON(t *testing.T) {
	must := func(file string) io.Reader {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		return f
	}

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
				err: nil,
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
				err: nil,
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
				err: nil,
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
		err := json.NewDecoder(must(tc.file)).Decode(&obtained)
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
