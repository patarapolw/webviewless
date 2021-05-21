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

// LocateChrome returns a path to Google Chrome binary
func LocateBrave() BrowserProtocol {
	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			// TODO: check on macOS
			"/usr/bin/brave",
		}
	case "windows":
		paths = []string{
			// TODO: check on Windows
		}
	default:
		paths = []string{
			"/usr/bin/brave",
			"/snap/bin/brave",
		}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		var c Connector = BraveConnector{
			Path: path,
		}

		return BrowserProtocol{
			Path:      path,
			Type:      "brave",
			Connector: &c,
		}
	}

	return BrowserProtocol{}
}

// BraveConnector connector for Brave browser
type BraveConnector struct {
	Path string
}

func (c BraveConnector) Type() string {
	return "brave"
}

func (c BraveConnector) Connect(url string, config *ConnectorConfig) chan bool {
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
		args = append(args, "--app="+url)
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

func (BraveConnector) OnDisconnect(cb func(), config *ConnectorConfig) chan bool {
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
