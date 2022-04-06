package main

import (
  "fmt"
  "github.com/mrspeiser/gocli/menus"
  "os"
  "encoding/json"
  "reflect"
)

func main() {

  type Configuration struct {
    Name    string `json:"NAME"`
    Version string `json:"VERSION"`
    Order   string `json:"ORDER"`
  }

  file, _ := os.Open("conf.json")
  defer file.Close()
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  err := decoder.Decode(&configuration)
  if err != nil {
    fmt.Println("error:", err)
  }

  configval := reflect.ValueOf(configuration)
  typeof := configval.Type()

  for i := 0; i < configval.NumField(); i++ {
    k := typeof.Field(i).Name
    v := configval.Field(i).Interface()
    fmt.Printf("%s: %s\n",k,v)
  }

  fmt.Println("\nGO CLI")
  fmt.Println("---------------------")

  runningstate := make(map[string]string)
  menu.ShowMenu(&runningstate)
}
