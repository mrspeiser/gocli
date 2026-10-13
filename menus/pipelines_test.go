package menus

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPipelineDelete(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pipelines.cli")

	names, err := LoadPipelines(path)
	if err != nil {
		t.Fatal(err)
	}
	if names != nil {
		t.Fatalf("missing file loaded %#v", names)
	}

	if err := SavePipelines(path, []string{"alpha", "beta", "gamma"}); err != nil {
		t.Fatal(err)
	}
	if err := deletePipelineAt(path, 2); err != nil {
		t.Fatal(err)
	}
	got, err := LoadPipelines(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"alpha", "gamma"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("pipeline %d = %q, want %q", i, got[i], want[i])
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Fatalf("pipelines.cli mode = %o, want 644", info.Mode().Perm())
	}
}
