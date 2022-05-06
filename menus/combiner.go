package menus

import (
  "strings"
)

func CollateMenuOptions(uniquekey string,
                        menuoptions []string,
                        cli *map[string]Menu){
  menuops := menuoptions
  combined := GetMenuMap(uniquekey,menuops)
  MainStateMenu := Menu{menuname:uniquekey}
  MainStateMenuOpts := make([]MenuOption,len(combined[uniquekey]))

  for i := 0; i < len(menuops); i++ {
    stringOpt := menuops[i]
    splitOpt := strings.Split(stringOpt, " ")
    MainStateMenuOpts[i] = MenuOption{
      accessor:splitOpt[0],
      value:splitOpt[1],
      uniquemenukey:splitOpt[2],
      action:splitOpt[3]}
  }
  MainStateMenu.menuoptions = MainStateMenuOpts
  (*cli)[uniquekey] = MainStateMenu
}

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

  stateops := []string{
    "q quit main.state.quit exit",
    "1 option1 main.state.option1 menu",
    "2 option2 main.state.option2 menu"}

  allmenuoptions := [][]string{mainopts,stateops}

  for i :=0; i<len(allmenuoptions); i++ {
    uniquekey := allmenukeys[i]
    menuoptions := allmenuoptions[i]
    CollateMenuOptions(uniquekey, menuoptions, &cli)
  }

  return cli
}
