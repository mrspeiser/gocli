package menus

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func LoadSequences(dir string) ([]Sequence, error) {
	data, err := os.ReadFile(filepath.Join(dir, "sequences.json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var sequences []Sequence
	if err := json.Unmarshal(data, &sequences); err != nil {
		return nil, err
	}
	return sequences, nil
}

func SaveSequences(dir string, sequences []Sequence) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if sequences == nil {
		sequences = []Sequence{}
	}
	data, err := json.MarshalIndent(sequences, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, "sequences.json"), data, 0644)
}

func NewSequence(rt *Runtime) {
	name, err := rt.prompt("name-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if name == "" {
		fmt.Println("name is required")
		return
	}
	pipelines, err := LoadPipelineConfigs(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if len(pipelines) == 0 {
		fmt.Println("no pipelines")
		return
	}
	fmt.Println("pipeline number, or an empty line to finish")
	var steps []string
	for {
		fmt.Print(formatPipelines(pipelines))
		input, err := rt.prompt("pipeline-> ")
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		if input == "" {
			break
		}
		n, err := strconv.Atoi(input)
		if err != nil || n < 1 || n > len(pipelines) {
			fmt.Println("invalid choice")
			continue
		}
		steps = append(steps, pipelines[n-1].Name)
	}
	if len(steps) == 0 {
		fmt.Println("sequence needs a pipeline")
		return
	}
	sequences, err := LoadSequences(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	sequences = append(sequences, Sequence{Name: name, Pipelines: steps})
	if err := SaveSequences(rt.dataDir(), sequences); err != nil {
		fmt.Println("error:", err)
	}
}

func ShowSequences(rt *Runtime) {
	sequences, err := LoadSequences(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Print(formatSequences(sequences))
}

func RunSequence(rt *Runtime) {
	sequences, err := LoadSequences(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if len(sequences) == 0 {
		fmt.Println("no sequences")
		return
	}
	fmt.Print(formatSequences(sequences))
	n, ok := pickStepIndex(rt, "sequence-> ", len(sequences))
	if !ok {
		return
	}
	if err := ExecuteSequence(rt, sequences[n-1]); err != nil {
		fmt.Println("error:", err)
	}
}

func DeleteSequence(rt *Runtime) {
	sequences, err := LoadSequences(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if len(sequences) == 0 {
		fmt.Println("no sequences")
		return
	}
	fmt.Print(formatSequences(sequences))
	n, ok := pickStepIndex(rt, "sequence-> ", len(sequences))
	if !ok {
		return
	}
	updated := make([]Sequence, 0, len(sequences)-1)
	updated = append(updated, sequences[:n-1]...)
	updated = append(updated, sequences[n:]...)
	if err := SaveSequences(rt.dataDir(), updated); err != nil {
		fmt.Println("error:", err)
	}
}

// ExecuteSequence runs each named pipeline in order. A failure stops the rest.
func ExecuteSequence(rt *Runtime, sequence Sequence) error {
	if err := reloadSteps(rt); err != nil {
		return err
	}
	pipelines, err := LoadPipelineConfigs(rt.dataDir())
	if err != nil {
		return err
	}
	if rt.State == nil {
		rt.State = map[string]string{}
	}
	for _, name := range sequence.Pipelines {
		pipeline, ok := pipelineByName(pipelines, name)
		if !ok {
			return fmt.Errorf("pipeline %q not found", name)
		}
		rt.State["last.sequence"] = sequence.Name
		fmt.Printf("\nrunning %s\n", name)
		status, body, err := ExecutePipeline(rt, pipeline)
		if body != "" || err == nil {
			fmt.Printf("loaded %s\n%s\n%s\n", pipeline.Load, status, body)
		}
		if err != nil {
			return fmt.Errorf("pipeline %s: %w", name, err)
		}
	}
	return nil
}

func pipelineByName(pipelines []Pipeline, name string) (Pipeline, bool) {
	for _, pipeline := range pipelines {
		if pipeline.Name == name {
			return pipeline, true
		}
	}
	return Pipeline{}, false
}
