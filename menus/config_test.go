package menus

import "testing"

func TestNormalizeProtocol(t *testing.T) {
	got, err := normalizeProtocol(" HTTP ")
	if err != nil {
		t.Fatal(err)
	}
	if got != "http" {
		t.Fatalf("protocol = %q", got)
	}
	if _, err := normalizeProtocol(""); err != nil {
		t.Fatal(err)
	}
	if _, err := normalizeProtocol("mongodb"); err == nil {
		t.Fatal("expected unsupported protocol error")
	}
}

func TestSavedSourceWithoutProtocolStillExtracts(t *testing.T) {
	src := Source{Name: "weather1", URL: "https://example.com/data.json"}
	if err := src.validateExtract(); err != nil {
		t.Fatal(err)
	}
	if src.protocolName() != "http" || src.httpMethod() != "GET" {
		t.Fatalf("defaults = %s %s", src.protocolName(), src.httpMethod())
	}
}

func TestSavedDestinationWithoutProtocolStillLoads(t *testing.T) {
	dest := Destination{Name: "hook", URL: "https://example.com/hook"}
	if err := dest.validateLoad(); err != nil {
		t.Fatal(err)
	}
	if dest.protocolName() != "http" || dest.httpMethod() != "POST" || dest.contentType() != "text/plain" {
		t.Fatalf("defaults = %s %s %s", dest.protocolName(), dest.httpMethod(), dest.contentType())
	}
}
