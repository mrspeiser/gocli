package menus

import (
  "strings"
)

func collate(){
  //allmenus := [][]string
  mainopts := []string{
    "q quit main.quit",
    "1 state main.state",
    "2 pipelines main.pipeline",
    "3 sources main.sources",
    "4 destinations main.destinations",
    "5 data transforms main.transforms",
    "6 sequences main.sequences"}
  MainMenu := Menu{}
  MainMenuOpts := make([]MenuOption,len(mainopts))

  for i := 0; i < len(mainopts); i++ {
    stringOpt := mainopts[i]
    splitOpt := strings.Split(stringOpt, " ")
    MainMenuOpts[i] = MenuOption{accessor:splitOpt[0],value:splitOpt[1],uniquemenukey:splitOpt[2]}
  }
  MainMenu.menuoptions = MainMenuOpts
  PrintMenuOptions(MainMenu)
}
