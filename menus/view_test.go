package menus

import "testing"

func TestFormatSaved(t *testing.T) {
	got := formatSaved(
		map[string]string{"b": "2", "a": "1"},
		[]Source{{Name: "docs", URL: "https://example.com"}},
		[]Pipeline{{Name: "nightly", Extract: "docs", Transform: "shout", Load: "hook"}},
		[]Destination{{Name: "hook", URL: "https://example.com/hook"}},
		[]Transform{{Name: "shout", Op: "upper"}},
	)
	want := "" +
		"\nstate:\n" +
		"a: 1\n" +
		"b: 2\n" +
		"\nextract:\n" +
		"1) docs http GET https://example.com\n" +
		"\npipelines:\n" +
		"1) nightly docs shout hook\n" +
		"\nload:\n" +
		"1) hook http POST https://example.com/hook text/plain\n" +
		"\ntransform:\n" +
		"1) shout upper\n"
	if got != want {
		t.Fatalf("saved view =\n%s\nwant\n%s", got, want)
	}
}

func TestFormatSavedEmpty(t *testing.T) {
	got := formatSaved(nil, nil, nil, nil, nil)
	want := "" +
		"\nstate:\n" +
		"state is empty\n" +
		"\nextract:\n" +
		"no extracts\n" +
		"\npipelines:\n" +
		"no pipelines\n" +
		"\nload:\n" +
		"no loads\n" +
		"\ntransform:\n" +
		"no transforms\n"
	if got != want {
		t.Fatalf("empty view =\n%s\nwant\n%s", got, want)
	}
}
