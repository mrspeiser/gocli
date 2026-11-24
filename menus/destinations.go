package menus

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func NewDestination(rt *Runtime) {
	name, err := rt.prompt("name: ")
	if err != nil {
		PrintError(fmt.Sprintf("Error: %v", err))
		return
	}
	if name == "" {
		PrintWarning("Name is required")
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
	contentType, err := rt.prompt("content-type: ")
	if err != nil {
		PrintError(fmt.Sprintf("Error: %v", err))
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
		PrintError(fmt.Sprintf("Failed to save destination: %v", err))
	} else {
		PrintSuccess(GetIcon("success") + "Destination created successfully!")
	}
}

func ShowDestinations(rt *Runtime) {
	destinations, err := LoadDestinations(rt.dataDir())
	if err != nil {
		PrintError(fmt.Sprintf("Failed to load destinations: %v", err))
		return
	}
	rt.Destinations = destinations
	if len(destinations) == 0 {
		PrintWarning("No destinations configured")
		return
	}
	fmt.Print(formatDestinations(destinations))
}

func DeleteDestination(rt *Runtime) {
	if len(rt.Destinations) == 0 {
		PrintWarning("No loads available")
		return
	}
	fmt.Print(formatDestinations(rt.Destinations))
	input, err := rt.prompt("load: ")
	if err != nil {
		PrintError(fmt.Sprintf("Error: %v", err))
		return
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(rt.Destinations) {
		PrintWarning("Invalid load number")
		return
	}
	deletedName := rt.Destinations[n-1].Name
	rt.Destinations = append(rt.Destinations[:n-1], rt.Destinations[n:]...)
	if err := SaveDestinations(rt.dataDir(), rt.Destinations); err != nil {
		PrintError(fmt.Sprintf("Failed to save: %v", err))
	} else {
		PrintSuccess(fmt.Sprintf("%s Deleted destination: %s", GetIcon("delete"), deletedName))
	}
}

func SendDestination(rt *Runtime) {
	body := ""
	if rt.State != nil {
		body = rt.State["last.body"]
	}
	if body == "" {
		PrintWarning("No body to send")
		return
	}
	dest, ok := pickDestination(rt)
	if !ok {
		return
	}
	if err := dest.validateLoad(); err != nil {
		PrintError(fmt.Sprintf("Validation error: %v", err))
		return
	}
	
	PrintInfo(fmt.Sprintf("%s Sending data to %s", GetIcon("destination"), dest.Name))
	
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
		PrintSection("Response")
		PrintInfo(fmt.Sprintf("Status: %s", status))
		fmt.Println(respBody)
		PrintSuccess(GetIcon("success") + "Data sent successfully!")
	}
	if err != nil {
		PrintError(fmt.Sprintf("Error: %v", err))
	}
}

// PostBody POSTs body to rawURL and keeps at most 4KB of the response.
func PostBody(rawURL, body, contentType string) (string, string, error) {
	bodyReader := strings.NewReader(body)
	contentLength := int64(len(body))
	
	resp, err := PostWithProgress(rawURL, bodyReader, contentLength, headerContentType(contentType), "Uploading data")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	
	limited, readErr := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	return resp.Status, string(limited), readErr
}

func pickDestination(rt *Runtime) (Destination, bool) {
	if len(rt.Destinations) == 0 {
		PrintWarning("No loads available")
		return Destination{}, false
	}
	fmt.Print(formatDestinations(rt.Destinations))
	input, err := rt.prompt("load: ")
	if err != nil {
		PrintError(fmt.Sprintf("Error: %v", err))
		return Destination{}, false
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(rt.Destinations) {
		PrintWarning("Invalid load number")
		return Destination{}, false
	}
	return rt.Destinations[n-1], true
}
