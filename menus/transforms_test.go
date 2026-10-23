package menus

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAndApplyTransform(t *testing.T) {
	op, args, err := parseTransformOp("trim")
	if err != nil {
		t.Fatal(err)
	}
	got, err := applyTransform(Transform{Op: op, Args: args}, "  hi  ")
	if err != nil {
		t.Fatal(err)
	}
	if got != "hi" {
		t.Fatalf("trim = %q", got)
	}

	op, args, err = parseTransformOp("upper")
	if err != nil {
		t.Fatal(err)
	}
	got, err = applyTransform(Transform{Op: op, Args: args}, "Ab")
	if err != nil {
		t.Fatal(err)
	}
	if got != "AB" {
		t.Fatalf("upper = %q", got)
	}

	op, args, err = parseTransformOp("replace foo bar baz")
	if err != nil {
		t.Fatal(err)
	}
	if len(args) != 2 || args[0] != "foo" || args[1] != "bar baz" {
		t.Fatalf("args = %#v", args)
	}
	got, err = applyTransform(Transform{Op: op, Args: args}, "foo bar foo")
	if err != nil {
		t.Fatal(err)
	}
	if got != "bar baz bar bar baz" {
		t.Fatalf("replace = %q", got)
	}

	if _, _, err := parseTransformOp("lower"); err == nil {
		t.Fatal("expected unknown op error")
	}
	if _, _, err := parseTransformOp("trim extra"); err == nil {
		t.Fatal("expected extra arg error")
	}
}

func TestTransformRoundTrip(t *testing.T) {
	dir := t.TempDir()
	got, err := LoadTransforms(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("missing file loaded %#v", got)
	}

	want := []Transform{{Name: "swap", Op: "replace", Args: []string{"foo", "bar baz"}}}
	if err := SaveTransforms(dir, want); err != nil {
		t.Fatal(err)
	}
	got, err = LoadTransforms(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != want[0].Name || got[0].Op != want[0].Op || len(got[0].Args) != 2 || got[0].Args[0] != "foo" || got[0].Args[1] != "bar baz" {
		t.Fatalf("transforms = %+v", got)
	}
	info, err := os.Stat(filepath.Join(dir, "transforms.json"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Fatalf("transforms.json mode = %o, want 644", info.Mode().Perm())
	}
}
