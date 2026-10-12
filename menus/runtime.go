package menus

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultDataDir = "data"

// Source is a named HTTP endpoint stored in data/sources.json.
type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Runtime is the in-memory session passed to menu actions.
type Runtime struct {
	State   map[string]string
	Sources []Source
	DataDir string
	in      *bufio.Reader
}

func (rt *Runtime) dataDir() string {
	if rt.DataDir == "" {
		return defaultDataDir
	}
	return rt.DataDir
}

func (rt *Runtime) readLine() (string, error) {
	if rt.in == nil {
		rt.in = bufio.NewReader(os.Stdin)
	}
	input, err := rt.in.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(input, "\r\n"), nil
}

func (rt *Runtime) prompt(label string) (string, error) {
	fmt.Print(label)
	return rt.readLine()
}

func LoadState(dir string) (map[string]string, error) {
	data, err := os.ReadFile(filepath.Join(dir, "state.json"))
	if err != nil {
		return nil, err
	}
	state := map[string]string{}
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return state, nil
}

func SaveState(dir string, state map[string]string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if state == nil {
		state = map[string]string{}
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, "state.json"), data, 0644)
}

func LoadSources(dir string) ([]Source, error) {
	data, err := os.ReadFile(filepath.Join(dir, "sources.json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var sources []Source
	if err := json.Unmarshal(data, &sources); err != nil {
		return nil, err
	}
	return sources, nil
}

func SaveSources(dir string, sources []Source) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if sources == nil {
		sources = []Source{}
	}
	data, err := json.MarshalIndent(sources, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, "sources.json"), data, 0644)
}
