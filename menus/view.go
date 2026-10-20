package menus

import (
	"fmt"
	"sort"
	"strings"
)

func ShowState(rt *Runtime) {
	sources, pipelines, ok := loadSaved(rt)
	if !ok {
		return
	}
	fmt.Print(formatSaved(rt.State, sources, pipelines))
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

func loadSaved(rt *Runtime) ([]Source, []string, bool) {
	sources, err := LoadSources(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, false
	}
	rt.Sources = sources
	pipelines, err := LoadPipelines(rt.pipelinesFile())
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, false
	}
	return sources, pipelines, true
}

func formatSaved(state map[string]string, sources []Source, pipelines []string) string {
	var b strings.Builder
	b.WriteString("\nstate:\n")
	b.WriteString(formatState(state))
	b.WriteString("\nsources:\n")
	b.WriteString(formatSources(sources))
	b.WriteString("\npipelines:\n")
	b.WriteString(formatPipelines(pipelines))
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
