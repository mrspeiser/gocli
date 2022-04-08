package menus

import (
  "strings"
)

func collate() map[string]Menu {
  allmenukeys := []string{
  "main",
  "main.state",
  "main.pipeline",
  "main.sources",
  "main.destinations",
  "main.transforms",
  "main.sequences"}

  cli := GetCliMap(allmenukeys)

  mainopts := []string{
    "q quit main.quit exit",
    "1 state main.state menu",
    "2 pipelines main.pipeline menu",
    "3 sources main.sources menu",
    "4 destinations main.destinations menu",
    "5 transforms main.transforms menu",
    "6 sequences main.sequences menu"}

  combined := GetMenuMap("main",mainopts)
  MainMenu := Menu{menuname:"main"}
  MainMenuOpts := make([]MenuOption,len(combined["main"]))

  for i := 0; i < len(mainopts); i++ {
    stringOpt := mainopts[i]
    splitOpt := strings.Split(stringOpt, " ")
    MainMenuOpts[i] = MenuOption{
      accessor:splitOpt[0],
      value:splitOpt[1],
      uniquemenukey:splitOpt[2],
      action:splitOpt[3]}
  }
  MainMenu.menuoptions = MainMenuOpts
  cli["main"] = MainMenu


  stateops := []string{
    "q quit main.state.quit exit",
    "1 option1 main.state.option1 menu",
    "2 option2 main.state.option2 menu"}

  combined2 := GetMenuMap("main.state",stateops)
  MainStateMenu := Menu{menuname:"main.state"}
  MainStateMenuOpts := make([]MenuOption,len(combined2["main.state"]))

  for i := 0; i < len(stateops); i++ {
    stringOpt := stateops[i]
    splitOpt := strings.Split(stringOpt, " ")
    MainStateMenuOpts[i] = MenuOption{
      accessor:splitOpt[0],
      value:splitOpt[1],
      uniquemenukey:splitOpt[2],
      action:splitOpt[3]}
  }
  MainStateMenu.menuoptions = MainStateMenuOpts
  cli["main.state"] = MainStateMenu

  return cli
}
