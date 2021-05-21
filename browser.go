package webviewless

import (
	"net"
	"os/exec"
	"runtime"
	"strconv"
)

// BrowserProtocol holds browser specific protocols
type BrowserProtocol struct {
	// Path to the executable
	Path string
	// Surface type
	Type string
	// Debugging connector
	Connector *Connector
}

// LocateBrowser returns a path to OS specific browser binary
func LocateBrowser(preferredBrowser string) BrowserProtocol {
	p := BrowserProtocol{}

	switch preferredBrowser {
	case "chrome":
		p = LocateChrome()
	case "edge":
		p = LocateEdge()
	case "brave":
		p = LocateBrave()
	// case "firefox":
	// 	p = LocateFirefox()
	default:
		switch runtime.GOOS {
		case "darwin":
		case "windows":
			p = LocateEdge()
		default:
			// p = LocateFirefox()
			// ! In Firefox, remote debugging port needs to be enabled explicitly.
		}
	}

	if p.Path == "" {
		p = LocateChrome()
	}

	if p.Path == "" {
		p = LocateBrave()
	}

	if p.Path == "" {
		p = LocateEdge()
	}

	// if p.Path == "" {
	// 	p = LocateFirefox()
	// }

	return p
}

type ConnectorConfig struct {
	Port        int
	Width       int
	Height      int
	IsMaximized bool

	// Internal
	Cmd *exec.Cmd
}

// abstract struct from universal connector
type Connector interface {
	// Type of the browser - 'chrome' | 'edge' | 'firefox'
	Type() string
	// Connect connects a URL to the browser
	Connect(url string, config *ConnectorConfig) chan bool
	// OnDisconnect calls callback function on disconnect
	OnDisconnect(cb func(), config *ConnectorConfig) chan bool
}

// FindRandomPort finds a random port
func FindRandomPort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	err = listener.Close()
	if err != nil {
		panic(err)
	}

	p, err := strconv.Atoi(port)
	if err != nil || p == 0 {
		panic(err)
	}

	return p
}
