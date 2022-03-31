package main

import (
  "fmt"
  "github.com/mrspeiser/gocli/menus"
)

func main() {

  //reader := bufio.NewReader(os.Stdin)
  fmt.Println("GO CLI")
  fmt.Println("---------------------")
  o1 := make(map[string]string)
  menu.Menu(&o1)
}
