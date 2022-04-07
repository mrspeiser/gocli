package menus

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
  action        string
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
returns a map of the menu name and un-collated menu options
   get the Menu by accessing the unique menu key
*/
func GetMenuMap(name string, menuopts []string) map[string][]string {
  menumap := make(map[string][]string)
  menumap[name] = menuopts
  return menumap
}

/* 
returns a map of all cli menu items
   get the Menu by accessing the unique menu key
*/
func GetCliMap(keys []string) map[string]Menu{
  cli := make(map[string]Menu)
  for i:=0;i<len(keys);i++{
    cli[keys[i]] = Menu{menuname:keys[i]}
  }
  return cli
}

