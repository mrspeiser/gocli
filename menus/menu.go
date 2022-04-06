package menu

import (
  "bufio"
  "os"
  "fmt"
  "strings"
)

type MenuOptions struct {
  accessors []string
  menuitems []string
  options2 []string
  actions   []string
}

func (m *MenuOptions) fill_defaults(){
  mainopts := []string{
    "q) quit",
    "1) state",
    "2) pipelines ",
    "3) sources",
    "4) destinations",
    "5) data transforms",
    "6) sequences"}

  stateops := []string{
    "b) back",
    "1) runtime",
    "2) saved ",
    "3) save current",
    "4) load saved",
    "5) add",
    "6) remove"}

  m.accessors = mainopts
  m.menuitems = mainopts
  m.options2 = stateops
  m.actions = mainopts
}

func PrintOptions(options []string){

  i := 0;
  for i < 7 {
    fmt.Println(options[i])
    i = i+1
  }
}
func Options(){
  // need to move this to configuration file eventually so
  // the user can change
  // how to access the index and order of the menu values
  // but not the value of "quit","state","pipelines","sources"
  mainoptions := MenuOptions{}
  mainoptions.fill_defaults()
  PrintOptions(mainoptions.menuitems)
}

func Menu(o1 *map[string]string){
  reader := bufio.NewReader(os.Stdin)
  state := *o1
  Options()

  for {
    fmt.Print("-> ")
    text, _ := reader.ReadString('\n')
    // convert CRLF to LF
    text = strings.Replace(text, "\n", "", -1)

    if strings.Compare("q", text) == 0 {
      break
    }

    if strings.Compare("hi", text) == 0 {
      fmt.Println("hello, Yourself")
    }

    if strings.Compare("rd", text) == 0 {
      fmt.Println("")
      key, _ := reader.ReadString('\n')
      val, _ := reader.ReadString('\n')
      key = strings.Replace(key, "\n", "", -1)
      val = strings.Replace(val, "\n", "", -1)
      (*o1)[key] = val
    }

    if strings.Compare("o1", text) == 0 {
      fmt.Println("\nMap o1:\n")
      for k,v := range state {
        fmt.Printf("%s: %s\n",k,v)
      }
      fmt.Println("\n")
    }
  }
}
