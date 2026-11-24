package menus

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func AccessMenu(cli map[string]Menu, actions map[string]func(*Runtime), rt *Runtime, currentmenu string) {
	for {
		PrintMenuOptions(cli[currentmenu])
		fmt.Print(StyledPrompt("Enter option: "))
		input, err := rt.readLine()
		if err != nil {
			PrintError(fmt.Sprintf("Error reading input: %v", err))
			break
		}

		if strings.Compare("q", input) == 0 {
			PrintInfo(GetIcon("quit") + "Exiting menu...")
			break
		}

		option, validoption := CheckOption(cli, input, currentmenu)

		if !validoption {
			PrintWarning("Invalid option, please try again")
			continue
		}

		if option != nil && option.action == "action" {
			fn := actions[option.uniquemenukey]
			if fn == nil {
				PrintError(fmt.Sprintf("No action found for %s", option.uniquemenukey))
			} else {
				PrintSection(option.value)
				fn(rt)
			}
		}

		if validoption && option.action == "menu" {
			nextaccess := option.uniquemenukey
			menu := cli[nextaccess]
			if menu.menuoptions != nil {
				AccessMenu(cli, actions, rt, nextaccess)
				PrintMenuOptions(cli[currentmenu])
			} else {
				PrintError(fmt.Sprintf("Menu %s has no options", menu.name))
			}
		}
	}
}

/*
Takes the input from the command line
Checks to see if the current menu's options
has a matching accessor, if it does
it returns the pointer to the menu option
*/
func CheckOption(cli map[string]Menu, input string, currentmenu string) (*MenuOption, bool) {
	menu := cli[currentmenu]
	if menu.menuoptions != nil {
		for i := 0; i < len(menu.menuoptions); i++ {
			if menu.menuoptions[i].accessor == input {
				menuoption := menu.menuoptions[i]
				return &menuoption, true
			}
		}
	}
	return nil, false
}

func ShowMenu(o1 *map[string]string) {
	rt := &Runtime{
		State:   *o1,
		DataDir: defaultDataDir,
	}
	
	PrintSection("Initializing Application")
	
	if state, err := LoadState(rt.DataDir); err == nil {
		rt.State = state
		*o1 = state
		PrintSuccess("State loaded successfully")
	} else if !errors.Is(err, os.ErrNotExist) {
		PrintError(fmt.Sprintf("Failed to load state: %v", err))
	}
	
	sources, err := LoadSources(rt.DataDir)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to load sources: %v", err))
	} else {
		rt.Sources = sources
		PrintSuccess(fmt.Sprintf("Loaded %d source(s)", len(sources)))
	}
	
	destinations, err := LoadDestinations(rt.DataDir)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to load destinations: %v", err))
	} else {
		rt.Destinations = destinations
		PrintSuccess(fmt.Sprintf("Loaded %d destination(s)", len(destinations)))
	}
	
	transforms, err := LoadTransforms(rt.DataDir)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to load transforms: %v", err))
	} else {
		rt.Transforms = transforms
		PrintSuccess(fmt.Sprintf("Loaded %d transform(s)", len(transforms)))
	}
	
	cli, actions := collate()
	AccessMenu(cli, actions, rt, "main")
}
