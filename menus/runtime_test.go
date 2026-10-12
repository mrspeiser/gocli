package menus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSourceRoundTrip(t *testing.T) {
	dir := t.TempDir()

	got, err := LoadSources(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("missing file loaded %#v", got)
	}

	want := []Source{
		{Name: "docs", URL: "https://example.com/docs"},
		{Name: "local", URL: "http://127.0.0.1:8080"},
	}
	if err := SaveSources(dir, want); err != nil {
		t.Fatal(err)
	}
	got, err = LoadSources(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("source %d = %+v, want %+v", i, got[i], want[i])
		}
	}

	info, err := os.Stat(filepath.Join(dir, "sources.json"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Fatalf("sources.json mode = %o, want 644", info.Mode().Perm())
	}
}

func TestFetchSource(t *testing.T) {
	t.Run("stores status and body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "hello from source")
		}))
		defer srv.Close()

		dir := t.TempDir()
		rt := &Runtime{State: map[string]string{}, DataDir: dir}
		body, err := FetchSource(rt, Source{Name: "demo", URL: srv.URL})
		if err != nil {
			t.Fatal(err)
		}
		if body != "hello from source" {
			t.Fatalf("body = %q", body)
		}
		if rt.State["last.url"] != srv.URL || rt.State["last.status"] != "200 OK" || rt.State["last.body"] != body {
			t.Fatalf("state = %#v", rt.State)
		}

		loaded, err := LoadState(dir)
		if err != nil {
			t.Fatal(err)
		}
		if loaded["last.body"] != body {
			t.Fatalf("saved body = %q", loaded["last.body"])
		}
	})

	t.Run("limits body to 4KB", func(t *testing.T) {
		payload := strings.Repeat("a", maxBodyBytes+1000)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, payload)
		}))
		defer srv.Close()

		rt := &Runtime{State: map[string]string{}, DataDir: t.TempDir()}
		body, err := FetchSource(rt, Source{Name: "big", URL: srv.URL})
		if err != nil {
			t.Fatal(err)
		}
		if len(body) != maxBodyBytes {
			t.Fatalf("len(body) = %d, want %d", len(body), maxBodyBytes)
		}
	})
}
