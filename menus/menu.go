package menus

import (
  "bufio"
  "os"
  "fmt"
  "strings"
  "github.com/mrspeiser/gocli/utilities"
)

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

func ShowMenu(o1 *map[string]string){
  reader := bufio.NewReader(os.Stdin)
  state := *o1
  mainopts := []string{
    "q) quit",
    "1) state",
    "2) pipelines ",
    "3) sources",
    "4) destinations",
    "5) data transforms",
    "6) sequences"}
  utils.PrintOptions(mainopts)

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
