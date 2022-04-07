package structures

/*
Menu Options struct holds:
  the accessor: The muatable value to access the next menu or action
  the unique menu key: The key on the MenuMap to get the next sub menu or previous menu
  the value: The immutable value that displays the option
*/
type MenuOption struct {
  uniquemenukey string
  accessor      string
  value         string
}

/*
Menu struct holds:
  The menuname: matches the uniquemenukey for the option selected
  the menuoptions: a list of MenuOptions
*/
type Menu struct {
  menuname    string
  menuoptions []MenuOption
}

/* 
returns a map of all cli menu items
   get the Menu by accessing the unique menu key
*/
func getmap() map[string]Menu{
  allmenus := make(map[string]Menu)
  return allmenus
}

