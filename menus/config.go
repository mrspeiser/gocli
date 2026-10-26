package menus

import (
	"fmt"
	"strings"
)

func askProtocol(rt *Runtime) (string, bool) {
	fmt.Println("protocols: http")
	raw, err := rt.prompt("protocol-> ")
	if err != nil {
		fmt.Println("error:", err)
		return "", false
	}
	protocol, err := normalizeProtocol(raw)
	if err != nil {
		fmt.Println("error:", err)
		return "", false
	}
	return protocol, true
}

func askHTTPURL(rt *Runtime) (string, bool) {
	rawURL, err := rt.prompt("url-> ")
	if err != nil {
		fmt.Println("error:", err)
		return "", false
	}
	if !validSourceURL(rawURL) {
		fmt.Println("url must be http or https")
		return "", false
	}
	return rawURL, true
}

func normalizeProtocol(raw string) (string, error) {
	protocol := strings.ToLower(strings.TrimSpace(raw))
	if protocol == "" {
		protocol = "http"
	}
	if protocol != "http" {
		return "", fmt.Errorf("only http is supported")
	}
	return protocol, nil
}

func (s Source) protocolName() string {
	if s.Protocol == "" {
		return "http"
	}
	return strings.ToLower(s.Protocol)
}

func (s Source) httpMethod() string {
	if s.Method == "" {
		return "GET"
	}
	return strings.ToUpper(s.Method)
}

func (s Source) validateExtract() error {
	if s.protocolName() != "http" {
		return fmt.Errorf("protocol %s is not supported", s.protocolName())
	}
	if s.httpMethod() != "GET" {
		return fmt.Errorf("extract method must be GET")
	}
	if !validSourceURL(s.URL) {
		return fmt.Errorf("url must be http or https")
	}
	return nil
}

func (d Destination) protocolName() string {
	if d.Protocol == "" {
		return "http"
	}
	return strings.ToLower(d.Protocol)
}

func (d Destination) httpMethod() string {
	if d.Method == "" {
		return "POST"
	}
	return strings.ToUpper(d.Method)
}

func (d Destination) contentType() string {
	if strings.TrimSpace(d.ContentType) == "" {
		return "text/plain"
	}
	return d.ContentType
}

func (d Destination) validateLoad() error {
	if d.protocolName() != "http" {
		return fmt.Errorf("protocol %s is not supported", d.protocolName())
	}
	if d.httpMethod() != "POST" {
		return fmt.Errorf("load method must be POST")
	}
	if !validSourceURL(d.URL) {
		return fmt.Errorf("url must be http or https")
	}
	return nil
}

func headerContentType(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "text/plain") {
		return "text/plain; charset=utf-8"
	}
	return raw
}
