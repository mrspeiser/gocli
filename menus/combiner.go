package combiner

import (
  "bufio"
  "os"
  "fmt"
  "strings"
  "github.com/mrspeiser/gocli/menus/menu"
)

func (m *MenuOption) fill_defaults(){
  mainopts := []string{
    "q) quit",
    "1) state",
    "2) pipelines ",
    "3) sources",
    "4) destinations",
    "5) data transforms",
    "6) sequences"}

  stateops := []string{
    "b) back",
    "1) runtime",
    "2) saved ",
    "3) save current",
    "4) load saved",
    "5) add",
    "6) remove"}
}
