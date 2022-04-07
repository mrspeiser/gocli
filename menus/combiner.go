package menus

import (
  "strings"
)

func collate(){
  allmenukeys := []string{
  "main",
  "main.state",
  "main.pipeline",
  "main.sources",
  "main.destinations",
  "main.transforms",
  "main.sequences"}

  menus := GetCliMap(allmenukeys)

  mainopts := []string{
    "q quit main.quit exit",
    "1 state main.state menu",
    "2 pipelines main.pipeline menu",
    "3 sources main.sources menu",
    "4 destinations main.destinations menu",
    "5 data transforms main.transforms menu",
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
  menus["main"] = MainMenu
  PrintMenuOptions(menus["main"])
}
