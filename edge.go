package webviewless

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

// LocateEdge returns a path to Microsoft Edge binary
// TODO: Update if MS Edge is installable on non-Windows.
func LocateEdge() BrowserProtocol {
	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			"/usr/bin/microsoft-edge",
			"/usr/bin/microsoft-edge-beta",
			"/usr/bin/microsoft-edge-dev",
		}
	case "windows":
		paths = []string{
			os.Getenv("ProgramFiles") + "/Microsoft/Edge/Application/msedge.exe",
			os.Getenv("ProgramFiles(x86)") + "/Microsoft/Edge/Application/msedge.exe",
		}
	default:
		paths = []string{
			"/usr/bin/microsoft-edge",
			"/usr/bin/microsoft-edge-beta",
			"/usr/bin/microsoft-edge-dev",
		}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		var c Connector = ChromeConnector{
			Path: path,
		}

		return BrowserProtocol{
			Path:      path,
			Type:      "edge",
			Connector: &c,
		}
	}
	return BrowserProtocol{}
}

// EdgeConnector connector for Microsoft Edge
type EdgeConnector struct {
	Path string
}

func (c EdgeConnector) Type() string {
	return "edge"
}

func (c EdgeConnector) Connect(url string, config *ConnectorConfig) chan bool {
	if config.Port == 0 {
		config.Port = FindRandomPort()
	}

	out := make(chan bool)

	go func() {
		for {
			_, err := net.Listen("tcp", ":"+strconv.Itoa(config.Port))
			if err != nil {
				break
			}

			time.Sleep(1 * time.Second)
		}

		out <- true
	}()

	args := []string{}
	if url != "" {
		args = append(args, "--app", url)
	}
	args = append(args, "--remote-debugging-port="+strconv.Itoa(config.Port))

	if config.IsMaximized {
		args = append(args, "--start-maximized")
	} else {
		if config.Height == 0 {
			config.Height = 600
		}
		if config.Width == 0 {
			config.Width = 800
		}

		args = append(args, fmt.Sprintf("--window-size=%d,%d", config.Width, config.Height))
	}

	config.Cmd = exec.Command(c.Path, args...)
	if e := config.Cmd.Start(); e != nil {
		panic(e)
	}

	return out
}

func (EdgeConnector) OnDisconnect(cb func(), config *ConnectorConfig) chan bool {
	if config.Port == 0 {
		panic("Cannot wait for port 0")
	}

	out := make(chan bool)

	go func() {
		if e := config.Cmd.Wait(); e != nil {
			panic(e)
		}

		cb()

		out <- true
	}()

	return out
}
