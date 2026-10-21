package menus

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDestinationRoundTrip(t *testing.T) {
	dir := t.TempDir()

	got, err := LoadDestinations(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("missing file loaded %#v", got)
	}

	want := []Destination{{Name: "hook", URL: "https://example.com/hook"}}
	if err := SaveDestinations(dir, want); err != nil {
		t.Fatal(err)
	}
	got, err = LoadDestinations(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("destinations = %+v", got)
	}

	info, err := os.Stat(filepath.Join(dir, "destinations.json"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Fatalf("destinations.json mode = %o, want 644", info.Mode().Perm())
	}
}

func TestPostBody(t *testing.T) {
	var gotMethod, gotBody, gotType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotType = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, "accepted")
	}))
	defer srv.Close()

	status, body, err := PostBody(srv.URL, "payload")
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPost || gotBody != "payload" || !strings.HasPrefix(gotType, "text/plain") {
		t.Fatalf("method=%s body=%q type=%s", gotMethod, gotBody, gotType)
	}
	if status != "200 OK" || body != "accepted" {
		t.Fatalf("status=%s body=%q", status, body)
	}
}
