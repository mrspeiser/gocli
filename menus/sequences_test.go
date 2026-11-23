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

func TestSequenceRoundTrip(t *testing.T) {
	dir := t.TempDir()
	got, err := LoadSequences(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("missing file loaded %#v", got)
	}
	want := []Sequence{{Name: "morning", Pipelines: []string{"one", "two"}}}
	if err := SaveSequences(dir, want); err != nil {
		t.Fatal(err)
	}
	got, err = LoadSequences(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "morning" || len(got[0].Pipelines) != 2 || got[0].Pipelines[0] != "one" || got[0].Pipelines[1] != "two" {
		t.Fatalf("sequences = %+v", got)
	}
	info, err := os.Stat(filepath.Join(dir, "sequences.json"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Fatalf("sequences.json mode = %o, want 644", info.Mode().Perm())
	}
}

func TestExecuteSequenceRunsPipelinesInOrder(t *testing.T) {
	firstExtract := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello")
	}))
	defer firstExtract.Close()
	secondExtract := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "world")
	}))
	defer secondExtract.Close()

	var loaded []string
	firstLoad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		loaded = append(loaded, "first:"+string(b))
		fmt.Fprint(w, "ok")
	}))
	defer firstLoad.Close()
	secondLoad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		loaded = append(loaded, "second:"+string(b))
		fmt.Fprint(w, "ok")
	}))
	defer secondLoad.Close()

	dir := t.TempDir()
	pipelines := []Pipeline{
		{Name: "one", Extract: "docs", Transform: "shout", Load: "hook"},
		{Name: "two", Extract: "other", Transform: "shout", Load: "archive"},
	}
	if err := SavePipelineConfigs(dir, pipelines); err != nil {
		t.Fatal(err)
	}
	rt := &Runtime{
		State:   map[string]string{},
		DataDir: dir,
		Sources: []Source{
			{Name: "docs", Protocol: "http", URL: firstExtract.URL, Method: "GET"},
			{Name: "other", Protocol: "http", URL: secondExtract.URL, Method: "GET"},
		},
		Transforms: []Transform{{Name: "shout", Op: "upper"}},
		Destinations: []Destination{
			{Name: "hook", Protocol: "http", URL: firstLoad.URL, Method: "POST", ContentType: "text/plain"},
			{Name: "archive", Protocol: "http", URL: secondLoad.URL, Method: "POST", ContentType: "text/plain"},
		},
	}
	if err := SaveSources(dir, rt.Sources); err != nil {
		t.Fatal(err)
	}
	if err := SaveTransforms(dir, rt.Transforms); err != nil {
		t.Fatal(err)
	}
	if err := SaveDestinations(dir, rt.Destinations); err != nil {
		t.Fatal(err)
	}

	err := ExecuteSequence(rt, Sequence{Name: "morning", Pipelines: []string{"one", "two"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 2 || loaded[0] != "first:HELLO" || loaded[1] != "second:WORLD" {
		t.Fatalf("loaded = %#v", loaded)
	}
	if rt.State["last.sequence"] != "morning" || rt.State["last.body"] != "WORLD" {
		t.Fatalf("state = %#v", rt.State)
	}
}
