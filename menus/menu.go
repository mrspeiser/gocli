package menu

import (
  "bufio"
  "os"
  "fmt"
  "strings"
  "github.com/mrspeiser/gocli/menus/combiner"
  "github.com/mrspeiser/gocli/utilities/printing"
)

/*
Menu Options struct holds:
  the accessor: The muatable value to access the next menu or action
  the unique menu key: The key on the MenuMap to get the next sub menu or previous menu
  the value: The immutable value that displays the option
*/
type MenuOption struct {
  uniquemenukey string
  accessor      string
  value         string
}

/*
Menu struct holds:
  The menuname: matches the uniquemenukey for the option selected
  the menuoptions: a list of MenuOptions
*/
type Menu struct {
  menuname    string
  menuoptions []MenuOption
}

/* 
returns a map of all cli menu items
   get the Menu by accessing the unique menu key
*/
func getmap() map[string]Menu{
  allmenus := make(map[string]Menu)
  return allmenus
}

func AccessMenu(menu Menu){
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("-> ")
  input, _ := reader.ReadString('\n')
  // convert CRLF to LF
  input = strings.Replace(input, "\n", "", -1)

  for i := 0; i<len(menu.menuoptions); i++{
    //menuopt := reflect.ValueOf(menu.menuoptions[i])
    if strings.Compare(menu.menuoptions[i].value, input) == 0 {
      fmt.Println(menu.menuoptions[i])
    }
  }
}


//func Options(){
//  // need to move this to configuration file eventually so
//  // the user can change
//  // how to access the index and order of the menu values
//  // but not the value of "quit","state","pipelines","sources"
//  mainoptions := MenuOption{}
//  mainoptions.fill_defaults()
//  PrintOptions(mainoptions.menuitems)
//}

func ShowMenu(o1 *map[string]string){
  reader := bufio.NewReader(os.Stdin)
  state := *o1
  //Options()

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
