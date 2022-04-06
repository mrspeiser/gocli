package utils

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
