package webviewless

import (
	"log"
	"runtime"
	"testing"
)

func TestLocateChrome(t *testing.T) {
	loc := LocateChrome()
	if loc.Path == "" {
		t.Fatalf("OS: %s - Cannot find Google Chrome\n", runtime.GOOS)
	}

	t.Logf("OS: %s - Google Chrome is at %s", runtime.GOOS, loc.Path)

	if loc.Type != "chrome" {
		t.Fatal("Type must only be 'chrome'")
	}

	if (*loc.Connector).Type() != "chrome" {
		t.Fatal("Connector must only be 'chrome'")
	}
}

func TestChromeAppMode(t *testing.T) {
	loc := LocateChrome()
	cfg := ConnectorConfig{
		// Width:       300,
		// Height:      200,
		// Port:        9222,
		// IsMaximized: true,
	}

	c := *loc.Connector
	<-c.Connect("http://example.org", &cfg)
	<-c.OnDisconnect(func() {
		log.Println("disconnecting")
	}, &cfg)
}
