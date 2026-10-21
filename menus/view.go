package menus

import (
	"fmt"
	"sort"
	"strings"
)

func ShowState(rt *Runtime) {
	sources, pipelines, destinations, ok := loadSaved(rt)
	if !ok {
		return
	}
	fmt.Print(formatSaved(rt.State, sources, pipelines, destinations))
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
	pipelines, err := LoadPipelines(rt.pipelinesFile())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Print(formatPipelines(pipelines))
}

func loadSaved(rt *Runtime) ([]Source, []string, []Destination, bool) {
	sources, err := LoadSources(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, nil, false
	}
	rt.Sources = sources
	pipelines, err := LoadPipelines(rt.pipelinesFile())
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, nil, false
	}
	destinations, err := LoadDestinations(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, nil, false
	}
	rt.Destinations = destinations
	return sources, pipelines, destinations, true
}

func formatSaved(state map[string]string, sources []Source, pipelines []string, destinations []Destination) string {
	var b strings.Builder
	b.WriteString("\nstate:\n")
	b.WriteString(formatState(state))
	b.WriteString("\nsources:\n")
	b.WriteString(formatSources(sources))
	b.WriteString("\npipelines:\n")
	b.WriteString(formatPipelines(pipelines))
	b.WriteString("\ndestinations:\n")
	b.WriteString(formatDestinations(destinations))
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
		return "no sources\n"
	}
	var b strings.Builder
	for i, src := range sources {
		fmt.Fprintf(&b, "%d) %s %s\n", i+1, src.Name, src.URL)
	}
	return b.String()
}

func formatDestinations(destinations []Destination) string {
	if len(destinations) == 0 {
		return "no destinations\n"
	}
	var b strings.Builder
	for i, dest := range destinations {
		fmt.Fprintf(&b, "%d) %s %s\n", i+1, dest.Name, dest.URL)
	}
	return b.String()
}

func formatPipelines(pipelines []string) string {
	if len(pipelines) == 0 {
		return "no pipelines\n"
	}
	var b strings.Builder
	for i, name := range pipelines {
		fmt.Fprintf(&b, "%d) %s\n", i+1, name)
	}
	return b.String()
}
