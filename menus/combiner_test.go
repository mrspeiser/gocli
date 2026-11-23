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

func TestSequencesMenuCanRun(t *testing.T) {
	cli, actions := collate()
	opts := cli["main.sequences"].menuoptions
	if len(opts) != 5 || opts[3].accessor != "3" || opts[3].action != "action" {
		t.Fatalf("sequences options = %+v", opts)
	}
	if actions["main.sequences.runsequence"] == nil {
		t.Fatal("run sequence is not wired")
	}
}
