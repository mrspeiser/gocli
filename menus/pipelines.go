package menus

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

func LoadPipelineConfigs(dir string) ([]Pipeline, error) {
	data, err := os.ReadFile(filepath.Join(dir, "pipelines.json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var pipelines []Pipeline
	if err := json.Unmarshal(data, &pipelines); err != nil {
		return nil, err
	}
	return pipelines, nil
}

func SavePipelineConfigs(dir string, pipelines []Pipeline) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if pipelines == nil {
		pipelines = []Pipeline{}
	}
	data, err := json.MarshalIndent(pipelines, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, "pipelines.json"), data, 0644)
}

func NewPipeline(rt *Runtime) {
	if err := reloadSteps(rt); err != nil {
		fmt.Println("error:", err)
		return
	}
	name, err := rt.prompt("name-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if name == "" {
		fmt.Println("name is required")
		return
	}
	extract, ok := pickStep(rt, "extract", formatSources(rt.Sources), len(rt.Sources), sourceNames(rt.Sources))
	if !ok {
		return
	}
	transform, ok := pickStep(rt, "transform", formatTransforms(rt.Transforms), len(rt.Transforms), transformNames(rt.Transforms))
	if !ok {
		return
	}
	load, ok := pickStep(rt, "load", formatDestinations(rt.Destinations), len(rt.Destinations), destinationNames(rt.Destinations))
	if !ok {
		return
	}
	pipelines, err := LoadPipelineConfigs(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	pipelines = append(pipelines, Pipeline{Name: name, Extract: extract, Transform: transform, Load: load})
	if err := SavePipelineConfigs(rt.dataDir(), pipelines); err != nil {
		fmt.Println("error:", err)
	}
}

func RunPipeline(rt *Runtime) {
	pipelines, err := LoadPipelineConfigs(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if len(pipelines) == 0 {
		fmt.Println("no pipelines")
		return
	}
	fmt.Print(formatPipelines(pipelines))
	n, ok := pickStepIndex(rt, "pipeline-> ", len(pipelines))
	if !ok {
		return
	}
	if err := reloadSteps(rt); err != nil {
		fmt.Println("error:", err)
		return
	}
	status, body, err := ExecutePipeline(rt, pipelines[n-1])
	if body != "" || err == nil {
		fmt.Printf("\nloaded %s\n%s\n%s\n", pipelines[n-1].Load, status, body)
	}
	if err != nil {
		fmt.Println("error:", err)
	}
}

func DeletePipeline(rt *Runtime) {
	pipelines, err := LoadPipelineConfigs(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if len(pipelines) == 0 {
		fmt.Println("no pipelines")
		return
	}
	fmt.Print(formatPipelines(pipelines))
	n, ok := pickStepIndex(rt, "pipeline-> ", len(pipelines))
	if !ok {
		return
	}
	updated := make([]Pipeline, 0, len(pipelines)-1)
	updated = append(updated, pipelines[:n-1]...)
	updated = append(updated, pipelines[n:]...)
	if err := SavePipelineConfigs(rt.dataDir(), updated); err != nil {
		fmt.Println("error:", err)
	}
}

// ExecutePipeline reads the extract, transforms the body, and POSTs it to the load.
func ExecutePipeline(rt *Runtime, pipeline Pipeline) (string, string, error) {
	src, ok := sourceByName(rt.Sources, pipeline.Extract)
	if !ok {
		return "", "", fmt.Errorf("extract %q not found", pipeline.Extract)
	}
	transform, ok := transformByName(rt.Transforms, pipeline.Transform)
	if !ok {
		return "", "", fmt.Errorf("transform %q not found", pipeline.Transform)
	}
	dest, ok := destinationByName(rt.Destinations, pipeline.Load)
	if !ok {
		return "", "", fmt.Errorf("load %q not found", pipeline.Load)
	}
	body, err := FetchSource(rt, src)
	if err != nil {
		return "", "", err
	}
	updated, err := applyTransform(transform, body)
	if err != nil {
		return "", "", err
	}
	if rt.State == nil {
		rt.State = map[string]string{}
	}
	rt.State["last.body"] = updated
	rt.State["last.transform"] = transform.Name
	if err := dest.validateLoad(); err != nil {
		return "", "", err
	}
	status, respBody, err := PostBody(dest.URL, updated, dest.contentType())
	if status != "" {
		rt.State["last.destination"] = dest.URL
		rt.State["last.sent.status"] = status
		rt.State["last.sent.body"] = respBody
		rt.State["last.pipeline"] = pipeline.Name
		rt.State["last.load"] = dest.Name
		if saveErr := SaveState(rt.dataDir(), rt.State); err == nil {
			err = saveErr
		}
	}
	return status, respBody, err
}

func reloadSteps(rt *Runtime) error {
	sources, err := LoadSources(rt.dataDir())
	if err != nil {
		return err
	}
	transforms, err := LoadTransforms(rt.dataDir())
	if err != nil {
		return err
	}
	destinations, err := LoadDestinations(rt.dataDir())
	if err != nil {
		return err
	}
	rt.Sources = sources
	rt.Transforms = transforms
	rt.Destinations = destinations
	return nil
}

func pickStep(rt *Runtime, kind, printed string, count int, names []string) (string, bool) {
	if count == 0 {
		fmt.Printf("no %ss\n", kind)
		return "", false
	}
	fmt.Print(printed)
	n, ok := pickStepIndex(rt, kind+"-> ", count)
	if !ok {
		return "", false
	}
	return names[n-1], true
}

func pickStepIndex(rt *Runtime, label string, count int) (int, bool) {
	input, err := rt.prompt(label)
	if err != nil {
		fmt.Println("error:", err)
		return 0, false
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > count {
		fmt.Println("invalid choice")
		return 0, false
	}
	return n, true
}

func sourceNames(sources []Source) []string {
	names := make([]string, len(sources))
	for i, src := range sources {
		names[i] = src.Name
	}
	return names
}

func transformNames(transforms []Transform) []string {
	names := make([]string, len(transforms))
	for i, transform := range transforms {
		names[i] = transform.Name
	}
	return names
}

func destinationNames(destinations []Destination) []string {
	names := make([]string, len(destinations))
	for i, dest := range destinations {
		names[i] = dest.Name
	}
	return names
}

func sourceByName(sources []Source, name string) (Source, bool) {
	for _, src := range sources {
		if src.Name == name {
			return src, true
		}
	}
	return Source{}, false
}

func transformByName(transforms []Transform, name string) (Transform, bool) {
	for _, transform := range transforms {
		if transform.Name == name {
			return transform, true
		}
	}
	return Transform{}, false
}

func destinationByName(destinations []Destination, name string) (Destination, bool) {
	for _, dest := range destinations {
		if dest.Name == name {
			return dest, true
		}
	}
	return Destination{}, false
}
