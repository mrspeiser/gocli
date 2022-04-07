package menus

import (
  "fmt"
)

func PrintOptions(options []string){
  i := 0;
  for i < 7 {
    fmt.Println(options[i])
    i = i+1
  }
}

func PrintMenuOptions(menu Menu){
  for i:=0; i<len(menu.menuoptions);i++{
    fmt.Printf("%s) %s\n",menu.menuoptions[i].accessor,menu.menuoptions[i].value)
  }
}
