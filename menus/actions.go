package menus

import (
	"os"
)

func WriteFile(filename string, s string) error {
	return os.WriteFile(filename, []byte(s), 0644)
}

func AppendFile(filename string, s string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(s + "\n")
	return err
}

func OpenFileReadOnly(filename string) (*os.File, error) {
	return os.Open(filename)
}
