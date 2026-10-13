package menus

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadPipelines(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	text := strings.TrimRight(string(data), "\r\n")
	if text == "" {
		return nil, nil
	}
	lines := strings.Split(text, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines, nil
}

func SavePipelines(filename string, names []string) error {
	var b strings.Builder
	for _, name := range names {
		b.WriteString(name)
		b.WriteByte('\n')
	}
	return os.WriteFile(filename, []byte(b.String()), 0644)
}

func WriteToPipelines(rt *Runtime) {
	input, err := rt.prompt("new pipeline-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	path := rt.pipelinesFile()
	names, err := LoadPipelines(path)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if err := SavePipelines(path, append(names, input)); err != nil {
		fmt.Println("error:", err)
	}
}

func DeletePipeline(rt *Runtime) {
	path := rt.pipelinesFile()
	names, err := LoadPipelines(path)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if len(names) == 0 {
		fmt.Println("no pipelines")
		return
	}
	for i, name := range names {
		fmt.Printf("%d) %s\n", i+1, name)
	}
	input, err := rt.prompt("pipeline-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(names) {
		fmt.Println("invalid pipeline")
		return
	}
	if err := deletePipelineAt(path, n); err != nil {
		fmt.Println("error:", err)
	}
}

func deletePipelineAt(filename string, n int) error {
	names, err := LoadPipelines(filename)
	if err != nil {
		return err
	}
	if n < 1 || n > len(names) {
		return fmt.Errorf("invalid pipeline")
	}
	updated := make([]string, 0, len(names)-1)
	updated = append(updated, names[:n-1]...)
	updated = append(updated, names[n:]...)
	return SavePipelines(filename, updated)
}
