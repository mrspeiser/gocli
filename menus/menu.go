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
		fmt.Print("-> ")
		input, err := rt.readLine()
		if err != nil {
			fmt.Println("error:", err)
			break
		}

		// exit loop if option is 'q'
		if strings.Compare("q", input) == 0 {
			break
		}

		// Print input to command line
		fmt.Println("\ninput:", input)

		// check the input to see if it has a matching option accessor
		option, validoption := CheckOption(cli, input, currentmenu)

		// Print input to command line
		fmt.Println("\noption:", option)

		if option != nil && option.action == "action" {
			fn := actions[option.uniquemenukey]
			if fn == nil {
				fmt.Println("no action for", option.uniquemenukey)
			} else {
				fn(rt)
			}
		}

		// if option is valid
		if validoption && option.action == "menu" {
			// get the next menu's key corresponding to the option
			nextaccess := option.uniquemenukey
			// retrieve menu from map using the options uniquemenukey
			menu := cli[nextaccess]
			// if menu's options are available
			if menu.menuoptions != nil {
				// recursively call AccessMenu
				AccessMenu(cli, actions, rt, nextaccess)
				PrintMenuOptions(cli[currentmenu])
			} else {
				fmt.Println("menu", menu, "has no options")
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
	if state, err := LoadState(rt.DataDir); err == nil {
		rt.State = state
		*o1 = state
	} else if !errors.Is(err, os.ErrNotExist) {
		fmt.Println("error:", err)
	}
	sources, err := LoadSources(rt.DataDir)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		rt.Sources = sources
	}
	destinations, err := LoadDestinations(rt.DataDir)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		rt.Destinations = destinations
	}
	cli, actions := collate()
	AccessMenu(cli, actions, rt, "main")
}
