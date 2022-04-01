package menu

import (
  "bufio"
  "os"
  "fmt"
  "strings"
)

func Options(){
  // need to move this to configuration file eventually so
  // the user can change
  // how to access the index and order of the menu values
  // but not the value of "quit","state","pipelines","sources"
  mainopts := [7]string{"q) quit","1) state","2) pipelines ","3) sources","4) destinations","5) data transforms","6) sequences"}
  i := 0;
  for i < 7 {
    fmt.Println(mainopts[i])
    i = i+1
  }
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
