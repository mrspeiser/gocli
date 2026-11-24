package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/mrspeiser/gocli/menus"
)

func main() {

	type Configuration struct {
		Name    string `json:"NAME"`
		Version string `json:"VERSION"`
		Order   string `json:"ORDER"`
	}

	file, err := os.Open("conf.json")
	if err != nil {
		fmt.Printf("error opening config file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Printf("error decoding config: %v\n", err)
		os.Exit(1)
	}

	configval := reflect.ValueOf(configuration)
	typeof := configval.Type()

	menus.PrintAppHeader(configuration.Name, configuration.Version)
	
	menus.PrintSection("Configuration")
	for i := 0; i < configval.NumField(); i++ {
		k := typeof.Field(i).Name
		v := configval.Field(i).Interface()
		menus.PrintInfo(fmt.Sprintf("%s: %v", k, v))
	}

	runningstate := make(map[string]string)
	menus.ShowMenu(&runningstate)
}
