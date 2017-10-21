package fs_test

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/kamilsk/stream/storage/fs"
)

var Update = flag.Bool("update", false, "update .golden files")

func TestEntity_WorkWithChildren(t *testing.T) {
	root := &fs.Entity{
		URN:  "root",
		Name: "My feeds",
	}
	child := &fs.Entity{URL: "https://rss.octolab.net/kamilsk/podcasts"}
	for _, tc := range []struct {
		name   string
		action func() bool
	}{
		{name: "add <nil> child", action: func() bool {
			return !root.AddChild(nil)
		}},
		{name: "add invalid child", action: func() bool {
			return !root.AddChild(&fs.Entity{})
		}},
		{name: "add child with duplicate URI", action: func() bool {
			root.AddChild(child)
			return !root.AddChild(&fs.Entity{URL: child.URL})
		}},
		{name: "remove <nil> child", action: func() bool {
			return !root.RemoveChild(nil)
		}},
		{name: "remove not presented child", action: func() bool {
			return !root.RemoveChild(&fs.Entity{URL: child.URL})
		}},
		{name: "remove child by empty URI", action: func() bool {
			return !root.RemoveChildByURI("")
		}},
		{name: "remove child by not presented URI", action: func() bool {
			return !root.RemoveChildByURI("unknown")
		}},
		{name: "successful removing child by URI", action: func() bool {
			return root.RemoveChildByURI(child.URI())
		}},
	} {
		if !tc.action() {
			t.Errorf("problem with test case %q", tc.name)
		}
	}
}

func TestEntity_MarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		entity   *fs.Entity
		expected struct {
			err    error
			golden string
		}
	}{
		{
			entity: &fs.Entity{URL: "https://rss.octolab.net/kamilsk/podcasts", Name: "end resource"},
			expected: struct {
				err    error
				golden string
			}{golden: "./fixtures/end_resource.golden"},
		},
		{
			entity: func() *fs.Entity {
				root := &fs.Entity{
					URN:  "multi_resource",
					Name: "multi resource",
				}
				root.AddChild(&fs.Entity{URL: "https://rss.octolab.net/kamilsk/podcasts", Name: "end resource"})
				root.AddChild(&fs.Entity{URL: "https://rss.octolab.net/kamilsk/releases", Name: "end resource"})
				root.AddChild(&fs.Entity{URL: "https://rss.octolab.net/octolab/releases", Name: "end resource"})
				return root
			}(),
			expected: struct {
				err    error
				golden string
			}{golden: "./fixtures/multi_resource.golden"},
		},
		{
			entity: func() *fs.Entity {
				root := &fs.Entity{
					URN:  "mixed_resource",
					Name: "mixed resource",
				}
				multi := &fs.Entity{
					URN:  "multi_resource",
					Name: "multi resource",
				}
				multi.AddChild(&fs.Entity{URL: "https://rss.octolab.net/kamilsk/podcasts", Name: "end resource"})
				multi.AddChild(&fs.Entity{URL: "https://rss.octolab.net/kamilsk/releases", Name: "end resource"})
				root.AddChild(multi)
				root.AddChild(&fs.Entity{URL: "https://rss.octolab.net/octolab/releases", Name: "end resource"})
				return root
			}(),
			expected: struct {
				err    error
				golden string
			}{golden: "./fixtures/mixed_resource.golden"},
		},
		{
			entity: func() *fs.Entity {
				// simple hack
				root := &fs.Entity{
					URN:  "invalid_resource",
					Name: "invalid resource",
				}
				child := &fs.Entity{URL: "https://rss.octolab.net/kamilsk/podcasts", Name: "end resource"}
				root.AddChild(child)
				child.Name, child.URL = "invalid_resource", ""
				return root
			}(),
			expected: struct {
				err    error
				golden string
			}{err: errors.New("cannot marshal invalid entity")},
		},
	} {
		data, err := json.Marshal(&tc.entity)
		switch {
		case tc.expected.err != nil && err != nil:
			if !strings.Contains(err.Error(), tc.expected.err.Error()) {
				t.Errorf("unexpected error. expected: %v, obtained: %v", tc.expected.err, err)
			}
		case tc.expected.err == nil && err != nil:
			fallthrough
		case tc.expected.err != nil && err == nil:
			t.Errorf("unexpected error. expected: %v, obtained: %v", tc.expected.err, err)
		}
		if tc.expected.golden == "" {
			continue
		}
		if *Update {
			json.NewEncoder(func(file string) io.Writer {
				f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				if err != nil {
					panic(err)
				}
				return f
			}(tc.expected.golden)).Encode(&tc.entity)
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

func TestEntity_UnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		file     string
		expected struct {
			err     error
			visitor func(fs.Entity) bool
		}
	}{
		{
			file: "./fixtures/end_resource.json",
			expected: struct {
				err     error
				visitor func(fs.Entity) bool
			}{
				visitor: func(e fs.Entity) bool {
					return e.IsValid() && e.URI() == "https://rss.octolab.net/kamilsk/podcasts" &&
						reflect.DeepEqual(e.Feeds(), []string{"https://rss.octolab.net/kamilsk/podcasts"})
				},
			},
		},
		{
			file: "./fixtures/multi_resource.json",
			expected: struct {
				err     error
				visitor func(fs.Entity) bool
			}{
				visitor: func(e fs.Entity) bool {
					return e.IsValid() && e.URI() == "multi_resource" && len(e.Children()) == 3 && func() bool {
						for i, src := range map[int]string{
							0: "https://rss.octolab.net/kamilsk/podcasts",
							1: "https://rss.octolab.net/kamilsk/releases",
							2: "https://rss.octolab.net/octolab/releases",
						} {
							if e.Children()[i].URI() != src {
								return false
							}
						}
						return true
					}() && reflect.DeepEqual(e.Feeds(), []string{
						"https://rss.octolab.net/kamilsk/podcasts",
						"https://rss.octolab.net/kamilsk/releases",
						"https://rss.octolab.net/octolab/releases",
					})
				},
			},
		},
		{
			file: "./fixtures/mixed_resource.json",
			expected: struct {
				err     error
				visitor func(fs.Entity) bool
			}{
				visitor: func(e fs.Entity) bool {
					return e.IsValid() && e.URI() == "mixed_resource" && len(e.Children()) == 2 &&
						e.Children()[1].URI() == "https://rss.octolab.net/octolab/releases" &&
						e.Children()[0].URI() == "mixed_resource/multi_resource" && func() bool {
						for i, src := range map[int]string{
							0: "https://rss.octolab.net/kamilsk/podcasts",
							1: "https://rss.octolab.net/kamilsk/releases",
						} {
							if e.Children()[0].Children()[i].URI() != src {
								return false
							}
						}
						return true
					}() && reflect.DeepEqual(e.Feeds(), []string{
						"https://rss.octolab.net/kamilsk/podcasts",
						"https://rss.octolab.net/kamilsk/releases",
						"https://rss.octolab.net/octolab/releases",
					})
				},
			},
		},
		{
			file: "./fixtures/invalid_resource.json",
			expected: struct {
				err     error
				visitor func(fs.Entity) bool
			}{err: errors.New("json: cannot unmarshal number into Go value of type []*fs.Entity")},
		},
	} {
		var obtained fs.Entity
		err := json.NewDecoder(func(file string) io.Reader {
			f, err := os.Open(file)
			if err != nil {
				panic(err)
			}
			return f
		}(tc.file)).Decode(&obtained)
		switch {
		case tc.expected.err != nil && err != nil:
			if !strings.Contains(err.Error(), tc.expected.err.Error()) {
				t.Errorf("unexpected error. expected: %v, obtained: %v", tc.expected.err, err)
			}
		case tc.expected.err == nil && err != nil:
			fallthrough
		case tc.expected.err != nil && err == nil:
			t.Errorf("unexpected error. expected: %v, obtained: %v", tc.expected.err, err)
		}
		if tc.expected.visitor == nil {
			continue
		}
		if !tc.expected.visitor(obtained) {
			t.Errorf("visitor failed at %s %q", obtained.URI(), obtained.Name)
		}
	}
}
