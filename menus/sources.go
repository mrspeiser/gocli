package menus

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const maxBodyBytes = 4096

func validSourceURL(raw string) bool {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	return parsed.Host != ""
}

func NewSource(rt *Runtime) {
	name, err := rt.prompt("name-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if name == "" {
		fmt.Println("name is required")
		return
	}
	rawURL, err := rt.prompt("url-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if !validSourceURL(rawURL) {
		fmt.Println("url must be http or https")
		return
	}
	rt.Sources = append(rt.Sources, Source{Name: name, URL: rawURL})
	if err := SaveSources(rt.dataDir(), rt.Sources); err != nil {
		fmt.Println("error:", err)
	}
}

func ReadSource(rt *Runtime) {
	src, ok := pickSource(rt)
	if !ok {
		return
	}
	body, err := FetchSource(rt, src)
	if body != "" || err == nil {
		fmt.Printf("\n%s\n%s\n", rt.State["last.status"], body)
	}
	if err != nil {
		fmt.Println("error:", err)
	}
}

func DeleteSource(rt *Runtime) {
	if len(rt.Sources) == 0 {
		fmt.Println("no sources")
		return
	}
	printSources(rt.Sources)
	input, err := rt.prompt("source-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(rt.Sources) {
		fmt.Println("invalid source")
		return
	}
	rt.Sources = append(rt.Sources[:n-1], rt.Sources[n:]...)
	if err := SaveSources(rt.dataDir(), rt.Sources); err != nil {
		fmt.Println("error:", err)
	}
}

func ShowState(rt *Runtime) {
	fmt.Println()
	if len(rt.State) == 0 {
		fmt.Println("state is empty")
		return
	}
	keys := make([]string, 0, len(rt.State))
	for k := range rt.State {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, rt.State[k])
	}
}

func SetState(rt *Runtime) {
	key, err := rt.prompt("key-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if key == "" {
		fmt.Println("key is required")
		return
	}
	val, err := rt.prompt("value-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if rt.State == nil {
		rt.State = map[string]string{}
	}
	rt.State[key] = val
	if err := SaveState(rt.dataDir(), rt.State); err != nil {
		fmt.Println("error:", err)
	}
}

// FetchSource GETs src.URL, keeps at most 4KB of the body, and stores
// last.url, last.status, and last.body on the runtime state.
func FetchSource(rt *Runtime, src Source) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(src.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	limited, readErr := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	body := string(limited)
	if rt.State == nil {
		rt.State = map[string]string{}
	}
	rt.State["last.url"] = src.URL
	rt.State["last.status"] = resp.Status
	rt.State["last.body"] = body
	if saveErr := SaveState(rt.dataDir(), rt.State); readErr == nil {
		readErr = saveErr
	}
	return body, readErr
}

func pickSource(rt *Runtime) (Source, bool) {
	if len(rt.Sources) == 0 {
		fmt.Println("no sources")
		return Source{}, false
	}
	printSources(rt.Sources)
	input, err := rt.prompt("source-> ")
	if err != nil {
		fmt.Println("error:", err)
		return Source{}, false
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(rt.Sources) {
		fmt.Println("invalid source")
		return Source{}, false
	}
	return rt.Sources[n-1], true
}

func printSources(sources []Source) {
	for i, src := range sources {
		fmt.Printf("%d) %s %s\n", i+1, src.Name, src.URL)
	}
}
