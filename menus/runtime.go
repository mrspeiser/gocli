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

const (
	defaultDataDir       = "data"
	defaultPipelinesFile = "pipelines.cli"
)

// Source is an extract configuration stored in data/sources.json.
type Source struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol,omitempty"`
	URL      string `json:"url,omitempty"`
	Method   string `json:"method,omitempty"`
}

// Destination is a load configuration stored in data/destinations.json.
type Destination struct {
	Name        string `json:"name"`
	Protocol    string `json:"protocol,omitempty"`
	URL         string `json:"url,omitempty"`
	Method      string `json:"method,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

// Pipeline is a saved extract, transform, and load, stored in data/pipelines.json.
type Pipeline struct {
	Name      string `json:"name"`
	Extract   string `json:"extract"`
	Transform string `json:"transform"`
	Load      string `json:"load"`
}

// Transform is a named text operation stored in data/transforms.json.
type Transform struct {
	Name string   `json:"name"`
	Op   string   `json:"op"`
	Args []string `json:"args,omitempty"`
}

// Runtime is the in-memory session passed to menu actions.
type Runtime struct {
	State         map[string]string
	Sources       []Source
	Destinations  []Destination
	Transforms    []Transform
	DataDir       string
	PipelinesFile string
	in            *bufio.Reader
}

func (rt *Runtime) dataDir() string {
	if rt.DataDir == "" {
		return defaultDataDir
	}
	return rt.DataDir
}

func (rt *Runtime) pipelinesFile() string {
	if rt.PipelinesFile == "" {
		return defaultPipelinesFile
	}
	return rt.PipelinesFile
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

func LoadDestinations(dir string) ([]Destination, error) {
	data, err := os.ReadFile(filepath.Join(dir, "destinations.json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var destinations []Destination
	if err := json.Unmarshal(data, &destinations); err != nil {
		return nil, err
	}
	return destinations, nil
}

func SaveDestinations(dir string, destinations []Destination) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if destinations == nil {
		destinations = []Destination{}
	}
	data, err := json.MarshalIndent(destinations, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, "destinations.json"), data, 0644)
}

func LoadTransforms(dir string) ([]Transform, error) {
	data, err := os.ReadFile(filepath.Join(dir, "transforms.json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var transforms []Transform
	if err := json.Unmarshal(data, &transforms); err != nil {
		return nil, err
	}
	return transforms, nil
}

func SaveTransforms(dir string, transforms []Transform) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if transforms == nil {
		transforms = []Transform{}
	}
	data, err := json.MarshalIndent(transforms, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, "transforms.json"), data, 0644)
}
