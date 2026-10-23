package menus

import "testing"

func TestCollateMenuOptionsParsesFields(t *testing.T) {
	cli := GetCliMap([]string{"main"})
	opts := []string{"1 state main.state menu"}
	CollateMenuOptions("main", opts, &cli)

	got := cli["main"].menuoptions
	if len(got) != 1 {
		t.Fatalf("len(menuoptions) = %d, want 1", len(got))
	}
	if got[0].accessor != "1" || got[0].value != "state" || got[0].uniquemenukey != "main.state" || got[0].action != "menu" {
		t.Fatalf("parsed option = %+v", got[0])
	}
}

func TestPlaceholderMenusHaveBack(t *testing.T) {
	cli, _ := collate()
	for _, key := range []string{"main.sequences"} {
		opts := cli[key].menuoptions
		if len(opts) != 1 || opts[0].accessor != "q" || opts[0].action != "exit" {
			t.Fatalf("%s options = %+v, want a single back option", key, opts)
		}
	}
}
