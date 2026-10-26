package menus

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewDestination(rt *Runtime) {
	name, err := rt.prompt("name-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if name == "" {
		fmt.Println("name is required")
		return
	}
	protocol, ok := askProtocol(rt)
	if !ok {
		return
	}
	rawURL, ok := askHTTPURL(rt)
	if !ok {
		return
	}
	contentType, err := rt.prompt("content-type-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	contentType = strings.TrimSpace(contentType)
	if contentType == "" {
		contentType = "text/plain"
	}
	rt.Destinations = append(rt.Destinations, Destination{
		Name:        name,
		Protocol:    protocol,
		URL:         rawURL,
		Method:      "POST",
		ContentType: contentType,
	})
	if err := SaveDestinations(rt.dataDir(), rt.Destinations); err != nil {
		fmt.Println("error:", err)
	}
}

func ShowDestinations(rt *Runtime) {
	destinations, err := LoadDestinations(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	rt.Destinations = destinations
	fmt.Print(formatDestinations(destinations))
}

func DeleteDestination(rt *Runtime) {
	if len(rt.Destinations) == 0 {
		fmt.Println("no loads")
		return
	}
	fmt.Print(formatDestinations(rt.Destinations))
	input, err := rt.prompt("load-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(rt.Destinations) {
		fmt.Println("invalid load")
		return
	}
	rt.Destinations = append(rt.Destinations[:n-1], rt.Destinations[n:]...)
	if err := SaveDestinations(rt.dataDir(), rt.Destinations); err != nil {
		fmt.Println("error:", err)
	}
}

func SendDestination(rt *Runtime) {
	body := ""
	if rt.State != nil {
		body = rt.State["last.body"]
	}
	if body == "" {
		fmt.Println("no body to send")
		return
	}
	dest, ok := pickDestination(rt)
	if !ok {
		return
	}
	if err := dest.validateLoad(); err != nil {
		fmt.Println("error:", err)
		return
	}
	status, respBody, err := PostBody(dest.URL, body, dest.contentType())
	if rt.State == nil {
		rt.State = map[string]string{}
	}
	if status != "" {
		rt.State["last.destination"] = dest.URL
		rt.State["last.sent.status"] = status
		rt.State["last.sent.body"] = respBody
		if saveErr := SaveState(rt.dataDir(), rt.State); err == nil {
			err = saveErr
		}
	}
	if respBody != "" || err == nil {
		fmt.Printf("\n%s\n%s\n", status, respBody)
	}
	if err != nil {
		fmt.Println("error:", err)
	}
}

// PostBody POSTs body to rawURL and keeps at most 4KB of the response.
func PostBody(rawURL, body, contentType string) (string, string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodPost, rawURL, strings.NewReader(body))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", headerContentType(contentType))
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	limited, readErr := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	return resp.Status, string(limited), readErr
}

func pickDestination(rt *Runtime) (Destination, bool) {
	if len(rt.Destinations) == 0 {
		fmt.Println("no loads")
		return Destination{}, false
	}
	fmt.Print(formatDestinations(rt.Destinations))
	input, err := rt.prompt("load-> ")
	if err != nil {
		fmt.Println("error:", err)
		return Destination{}, false
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(rt.Destinations) {
		fmt.Println("invalid load")
		return Destination{}, false
	}
	return rt.Destinations[n-1], true
}
