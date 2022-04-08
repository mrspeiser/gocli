package menus

import (
  "bufio"
  "os"
  "fmt"
  "strings"
)

func AccessMenu(cli map[string]Menu,o1 *map[string]string,currentmenu string){
  reader := bufio.NewReader(os.Stdin)
  // add for loop to continue if no menu
  fmt.Print("-> ")
  input, _ := reader.ReadString('\n')
  // convert CRLF to LF
  input = strings.Replace(input, "\n", "", -1)
  if strings.Compare("q", input) == 0 {
    return
  }
  fmt.Println("\n%s\n",input)
  result := CheckOption(cli,input,currentmenu)
  fmt.Printf(result)
  // check result to see if menu is available
  // if menu is available 
  // rereive menu from map 
  // recursively call AccessMenu
}

func CheckOption(cli map[string]Menu, input string, currentmenu string) string {
  menu := cli[currentmenu]
  if menu.menuoptions != nil {
    for i:=0; i<len(menu.menuoptions); i++ {
      if menu.menuoptions[i].accessor == input {
        fmt.Println("Found option: %s",menu.menuoptions[i])
        fmt.Println("next menu: ",menu.menuoptions[i].uniquemenukey)
        return menu.menuoptions[i].uniquemenukey
      }
    }
  }
  return "no menu available"
}

func ShowMenu(o1 *map[string]string){
  state := *o1
  cli := collate()
  AccessMenu(cli,&state,"main")
  //reader := bufio.NewReader(os.Stdin)
  //for {
  //  fmt.Print("-> ")
  //  text, _ := reader.ReadString('\n')
  //  // convert CRLF to LF
  //  text = strings.Replace(text, "\n", "", -1)

  //  if strings.Compare("q", text) == 0 {
  //    break
  //  }

  //  if strings.Compare("hi", text) == 0 {
  //    fmt.Println("hello, Yourself")
  //  }

  //  if strings.Compare("rd", text) == 0 {
  //    fmt.Println("")
  //    key, _ := reader.ReadString('\n')
  //    val, _ := reader.ReadString('\n')
  //    key = strings.Replace(key, "\n", "", -1)
  //    val = strings.Replace(val, "\n", "", -1)
  //    (*o1)[key] = val
  //  }

  //  if strings.Compare("o1", text) == 0 {
  //    fmt.Println("\nMap o1:\n")
  //    for k,v := range state {
  //      fmt.Printf("%s: %s\n",k,v)
  //    }
  //    fmt.Println("\n")
  //  }
  //}
}
