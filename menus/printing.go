package menus

import (
	"github.com/pterm/pterm"
)

func PrintOptions(options []string) {
	i := 0
	for i < 7 && i < len(options) {
		PrintInfo(options[i])
		i = i + 1
	}
}

func PrintMenuOptions(menu Menu) {
	pterm.Println()
	PrintMenuHeader(menu.name)
	
	tableData := pterm.TableData{
		{"Option", "Description"},
	}
	
	for i := 0; i < len(menu.menuoptions); i++ {
		opt := menu.menuoptions[i]
		icon := GetIcon(getIconForOption(opt.value))
		description := icon + opt.value
		tableData = append(tableData, []string{
			pterm.FgYellow.Sprint(opt.accessor),
			pterm.FgWhite.Sprint(description),
		})
	}
	
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()
}

func getIconForOption(value string) string {
	iconMap := map[string]string{
		"pipelines": "pipeline",
		"sources": "source",
		"destinations": "destination",
		"transforms": "transform",
		"sequences": "sequence",
		"state": "state",
		"new": "add",
		"show": "view",
		"delete": "delete",
		"run": "run",
		"send": "action",
		"read": "action",
		"quit": "quit",
		"back": "back",
		"set": "save",
	}
	
	for key, icon := range iconMap {
		if containsWord(value, key) {
			return icon
		}
	}
	
	return "menu"
}

func containsWord(s, word string) bool {
	for i := 0; i < len(s)-len(word)+1; i++ {
		match := true
		for j := 0; j < len(word); j++ {
			if toLower(s[i+j]) != toLower(word[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}
