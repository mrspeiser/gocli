package menus

import (
  "os"
  "fmt"
  "log"
)

func WriteFile(filename string, s string) int {
  fn := filename
  str := []byte(s)
  err := os.WriteFile(fn,str,0777)
  if err != nil {
    fmt.Printf("Error writing file: %s",err)
    return 1
  }
  return 0
}

func AppendFile(filename string, s string) {
  f, err := os.OpenFile(filename,
	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
  if err != nil {
    log.Println(err)
  }
  s = s+"\n"
  defer f.Close()
  if _, err := f.WriteString(s); err != nil {
    log.Println(err)
  }
}

func OpenFileReadOnly(filename string) (file *os.File) {
  fn := filename
  /*
    os.OpenFile(filename, safety_flag, permissions)
  */
  f,err := os.OpenFile(fn,os.O_RDONLY,755)
  if err != nil {
    fmt.Printf("Error opening file: %s",err)
    panic(err)
  }
  return f
}
