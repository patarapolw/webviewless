package webviewless

import (
	"log"
	"runtime"
	"testing"
)

func TestLocateEdge(t *testing.T) {
	loc := LocateEdge()
	if loc.Path == "" {
		t.Fatalf("OS: %s - Cannot find MS Edge\n", runtime.GOOS)
	}

	t.Logf("OS: %s - MS Edge is at %s", runtime.GOOS, loc.Path)

	if loc.Type != "edge" {
		t.Fatal("Type must only be 'edge'")
	}

	if (*loc.Connector).Type() != "edge" {
		t.Fatal("Connector must only be 'chrome'")
	}
}

func TestEdgeAppMode(t *testing.T) {
	loc := LocateEdge()
	cfg := ConnectorConfig{
		Width:  300,
		Height: 200,
		Port:   9222,
		// IsMaximized: true,
	}

	c := *loc.Connector
	<-c.Connect("http://example.org", &cfg)
	<-c.OnDisconnect(func() {
		log.Println("disconnecting")
	}, &cfg)
}
