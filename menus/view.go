package menus

import (
	"fmt"
	"sort"
	"strings"
)

func ShowState(rt *Runtime) {
	sources, pipelines, destinations, transforms, ok := loadSaved(rt)
	if !ok {
		return
	}
	fmt.Print(formatSaved(rt.State, sources, pipelines, destinations, transforms))
}

func ShowSources(rt *Runtime) {
	sources, err := LoadSources(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	rt.Sources = sources
	fmt.Print(formatSources(sources))
}

func ShowPipelines(rt *Runtime) {
	pipelines, err := LoadPipelineConfigs(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Print(formatPipelines(pipelines))
}

func loadSaved(rt *Runtime) ([]Source, []Pipeline, []Destination, []Transform, bool) {
	sources, err := LoadSources(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, nil, nil, false
	}
	rt.Sources = sources
	pipelines, err := LoadPipelineConfigs(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, nil, nil, false
	}
	destinations, err := LoadDestinations(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, nil, nil, false
	}
	rt.Destinations = destinations
	transforms, err := LoadTransforms(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, nil, nil, false
	}
	rt.Transforms = transforms
	return sources, pipelines, destinations, transforms, true
}

func formatSaved(state map[string]string, sources []Source, pipelines []Pipeline, destinations []Destination, transforms []Transform) string {
	var b strings.Builder
	b.WriteString("\nstate:\n")
	b.WriteString(formatState(state))
	b.WriteString("\nextract:\n")
	b.WriteString(formatSources(sources))
	b.WriteString("\npipelines:\n")
	b.WriteString(formatPipelines(pipelines))
	b.WriteString("\nload:\n")
	b.WriteString(formatDestinations(destinations))
	b.WriteString("\ntransform:\n")
	b.WriteString(formatTransforms(transforms))
	return b.String()
}

func formatState(state map[string]string) string {
	if len(state) == 0 {
		return "state is empty\n"
	}
	keys := make([]string, 0, len(state))
	for k := range state {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&b, "%s: %s\n", k, state[k])
	}
	return b.String()
}

func formatSources(sources []Source) string {
	if len(sources) == 0 {
		return "no extracts\n"
	}
	var b strings.Builder
	for i, src := range sources {
		fmt.Fprintf(&b, "%d) %s %s %s %s\n", i+1, src.Name, src.protocolName(), src.httpMethod(), src.URL)
	}
	return b.String()
}

func formatDestinations(destinations []Destination) string {
	if len(destinations) == 0 {
		return "no loads\n"
	}
	var b strings.Builder
	for i, dest := range destinations {
		fmt.Fprintf(&b, "%d) %s %s %s %s %s\n", i+1, dest.Name, dest.protocolName(), dest.httpMethod(), dest.URL, dest.contentType())
	}
	return b.String()
}

func formatTransforms(transforms []Transform) string {
	if len(transforms) == 0 {
		return "no transforms\n"
	}
	var b strings.Builder
	for i, transform := range transforms {
		fmt.Fprintf(&b, "%d) %s %s\n", i+1, transform.Name, formatTransformOp(transform))
	}
	return b.String()
}

func formatTransformOp(transform Transform) string {
	if transform.Op == "replace" && len(transform.Args) == 2 {
		return transform.Op + " " + transform.Args[0] + " " + transform.Args[1]
	}
	return transform.Op
}

func formatPipelines(pipelines []Pipeline) string {
	if len(pipelines) == 0 {
		return "no pipelines\n"
	}
	var b strings.Builder
	for i, pipeline := range pipelines {
		fmt.Fprintf(&b, "%d) %s %s %s %s\n", i+1, pipeline.Name, pipeline.Extract, pipeline.Transform, pipeline.Load)
	}
	return b.String()
}
