package webviewless

import (
	"log"
	"runtime"
	"testing"
)

func TestLocateBrave(t *testing.T) {
	loc := LocateBrave()
	if loc.Path == "" {
		t.Fatalf("OS: %s - Cannot find Brave browser\n", runtime.GOOS)
	}

	t.Logf("OS: %s - Brave browser is at %s", runtime.GOOS, loc.Path)

	if loc.Type != "brave" {
		t.Fatal("Type must only be 'brave'")
	}

	if (*loc.Connector).Type() != "chrome" {
		t.Fatal("Connector must only be 'chrome'")
	}
}

func TestBraveAppMode(t *testing.T) {
	loc := LocateBrave()
	cfg := ConnectorConfig{
		Width: 300,
		// Height: 200,
		Port: 9222,
		// IsMaximized: true,
	}

	c := *loc.Connector
	<-c.Connect("http://example.org", &cfg)
	<-c.OnDisconnect(func() {
		log.Println("disconnecting")
	}, &cfg)
}
