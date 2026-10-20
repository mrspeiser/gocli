package menus

import "testing"

func TestFormatSaved(t *testing.T) {
	got := formatSaved(
		map[string]string{"b": "2", "a": "1"},
		[]Source{{Name: "docs", URL: "https://example.com"}},
		[]string{"alpha"},
	)
	want := "" +
		"\nstate:\n" +
		"a: 1\n" +
		"b: 2\n" +
		"\nsources:\n" +
		"1) docs https://example.com\n" +
		"\npipelines:\n" +
		"1) alpha\n"
	if got != want {
		t.Fatalf("saved view =\n%s\nwant\n%s", got, want)
	}
}

func TestFormatSavedEmpty(t *testing.T) {
	got := formatSaved(nil, nil, nil)
	want := "" +
		"\nstate:\n" +
		"state is empty\n" +
		"\nsources:\n" +
		"no sources\n" +
		"\npipelines:\n" +
		"no pipelines\n"
	if got != want {
		t.Fatalf("empty view =\n%s\nwant\n%s", got, want)
	}
}
