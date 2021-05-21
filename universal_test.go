package webviewless

import (
	"runtime"
	"testing"
)

func TestLocateUniversal(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		loc := LocateBrowser("")
		t.Logf("OS: %s - MS Edge is at %s", runtime.GOOS, loc.Path)

		if loc.Type != "edge" {
			t.Fatal("Type must only be 'edge'")
		}

		if (*loc.Connector).Type() != "chrome" {
			t.Fatal("Connector must only be 'chrome'")
		}
	default:
		TestLocateChrome(t)
	}
}
