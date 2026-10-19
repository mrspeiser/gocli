package menus

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
		if path := rt.State["last.file"]; path != "" {
			fmt.Println("wrote", path)
		}
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

// FetchSource GETs src.URL, keeps at most 4KB of the body, stores
// last.url, last.status, and last.body, and writes the body to a file.
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
	delete(rt.State, "last.file")
	rt.State["last.url"] = src.URL
	rt.State["last.status"] = resp.Status
	rt.State["last.body"] = body
	if readErr == nil {
		path, writeErr := SaveSourceBody(rt.dataDir(), src.Name, body)
		if writeErr != nil {
			readErr = writeErr
		} else {
			rt.State["last.file"] = path
		}
	}
	if saveErr := SaveState(rt.dataDir(), rt.State); readErr == nil {
		readErr = saveErr
	}
	return body, readErr
}

// SaveSourceBody writes the response body to data/<name>.txt.
func SaveSourceBody(dir, name, body string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, sourceFileName(name))
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return "", err
	}
	return path, nil
}

func sourceFileName(name string) string {
	cleaned := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			return r
		default:
			return '_'
		}
	}, name)
	cleaned = strings.Trim(cleaned, "_")
	if cleaned == "" {
		return "source.txt"
	}
	return cleaned + ".txt"
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
