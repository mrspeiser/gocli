package menus

import (
  "bufio"
  "os"
  "fmt"
  "strings"
)

func AccessMenu(cli map[string]Menu,
                o1 *map[string]string,
                currentmenu string){
  reader := bufio.NewReader(os.Stdin)
  PrintMenuOptions(cli[currentmenu])
  for {
    fmt.Print("-> ")
    input, _ := reader.ReadString('\n')
    // convert CRLF to LF
    input = strings.Replace(input, "\n", "", -1)

    // exit loop if option is 'q'
    if strings.Compare("q", input) == 0 {
      break
    }

    // Print input to command line
    fmt.Println("\ninput: ",input,"\n")

    // check the input to see if it has a matching option accessor
    option,validoption := CheckOption(cli,input,currentmenu)

    // if option is valid
    if(validoption == true){
      // get the next menu's key corresponding to the option
      nextaccess := (*option).uniquemenukey
      // rereive menu from map using the options uniquemenukey 
      menu := cli[nextaccess]
      // if menu's options are available 
      if menu.menuoptions != nil {
        // recursively call AccessMenu
        AccessMenu(cli,o1,nextaccess)
        PrintMenuOptions(cli[currentmenu])
      // else tell user, the menu they selected is not setup
      } else {
        fmt.Println("menu ",menu," has no options\n")
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

func CheckOption(cli map[string]Menu, input string, currentmenu string) (*MenuOption,bool) {
  menu := cli[currentmenu]
  if menu.menuoptions != nil {
    for i:=0; i<len(menu.menuoptions); i++ {
      if menu.menuoptions[i].accessor == input {
        menuoption := menu.menuoptions[i]
        return &menuoption,true
      }
    }
  }
  return nil,false
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
