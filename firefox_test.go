package webviewless

import (
	"log"
	"runtime"
	"testing"
)

func TestLocateFirefox(t *testing.T) {
	loc := LocateFirefox()
	if loc.Path == "" {
		t.Fatalf("OS: %s - Cannot find Mozilla Firefox\n", runtime.GOOS)
	}

	t.Logf("OS: %s - Mozilla Firefox is at %s", runtime.GOOS, loc.Path)

	if loc.Type != "firefox" {
		t.Fatal("Type must only be 'firefox'")
	}

	if (*loc.Connector).Type() != "firefox" {
		t.Fatal("Connector must only be 'firefox'")
	}
}

func TestFirefoxAppMode(t *testing.T) {
	loc := LocateFirefox()
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
