package menus

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestPipelineConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	got, err := LoadPipelineConfigs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("missing file loaded %#v", got)
	}
	want := []Pipeline{{Name: "nightly", Extract: "docs", Transform: "shout", Load: "hook"}}
	if err := SavePipelineConfigs(dir, want); err != nil {
		t.Fatal(err)
	}
	got, err = LoadPipelineConfigs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("pipelines = %+v", got)
	}
	info, err := os.Stat(filepath.Join(dir, "pipelines.json"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Fatalf("pipelines.json mode = %o, want 644", info.Mode().Perm())
	}
}

func TestExecutePipelineLoadsTransformedBody(t *testing.T) {
	extract := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello")
	}))
	defer extract.Close()

	var gotBody, gotType string
	load := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		gotType = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, "stored")
	}))
	defer load.Close()

	dir := t.TempDir()
	rt := &Runtime{
		State:   map[string]string{},
		DataDir: dir,
		Sources: []Source{{
			Name:     "docs",
			Protocol: "http",
			URL:      extract.URL,
			Method:   "GET",
		}},
		Transforms: []Transform{{Name: "shout", Op: "upper"}},
		Destinations: []Destination{{
			Name:        "hook",
			Protocol:    "http",
			URL:         load.URL,
			Method:      "POST",
			ContentType: "text/plain",
		}},
	}
	status, body, err := ExecutePipeline(rt, Pipeline{
		Name:      "nightly",
		Extract:   "docs",
		Transform: "shout",
		Load:      "hook",
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotBody != "HELLO" {
		t.Fatalf("loaded body = %q", gotBody)
	}
	if status != "200 OK" || body != "stored" {
		t.Fatalf("status=%s body=%q", status, body)
	}
	if gotType == "" || rt.State["last.body"] != "HELLO" || rt.State["last.load"] != "hook" {
		t.Fatalf("type=%s state=%#v", gotType, rt.State)
	}
}
