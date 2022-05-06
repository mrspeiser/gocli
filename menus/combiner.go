package menus

import (
  "strings"
  "bufio"
  "os"
  "fmt"
)

func WriteToPipelines(){
  fmt.Print("new pipeline-> ")
  reader := bufio.NewReader(os.Stdin)
  input, _ := reader.ReadString('\n')
  input = strings.Replace(input, "\n", "", -1)
  AppendFile("pipelines.cli",input)
}

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

func collate() (map[string]Menu, map[string]func()) {

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
    "q back main.state.quit exit",
    "1 show-o1 main.state.option1 menu",
    "2 set-o1 main.state.option2 menu"}

  pipelineops := []string{
    "q back main.pipeline.quit exit",
    "1 new main.pipeline.newpipeline action",
    "2 delete main.pipeline.deletepipeline action"}

  sourcesops := []string{
    "q back main.sources.quit exit",
    "1 new main.sources.newsource action",
    "2 read main.sources.readsource action",
    "3 delete main.sources.deletesource action"}

  /* 
    IMPORTANT, when adding a new menu's options they need to be in the 
    correct order with the allmenukeys variable above, otherwise
    the menus will be mismatched
  */

  allmenuoptions := [][]string{mainopts,stateops,pipelineops,sourcesops}

  for i :=0; i<len(allmenuoptions); i++ {
    uniquekey := allmenukeys[i]
    menuoptions := allmenuoptions[i]
    CollateMenuOptions(uniquekey, menuoptions, &cli)
  }

  actions := make(map[string]func())
  actions["main.pipeline.newpipeline"] = WriteToPipelines

  return cli, actions
}

