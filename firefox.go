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

// LocateFirefox returns a path to Mozilla Firefox binary
func LocateFirefox() BrowserProtocol {
	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			// TODO: test on macOS
			"/usr/bin/firefox",
		}
	case "windows":
		paths = []string{
			// TODO: test on Windows
		}
	default:
		paths = []string{
			"/usr/bin/firefox",
		}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		var c Connector = FirefoxConnector{
			Path: path,
		}

		return BrowserProtocol{
			Path:      path,
			Type:      "firefox",
			Connector: &c,
		}
	}

	return BrowserProtocol{}
}

// FirefoxConnector connector for Firefox
type FirefoxConnector struct {
	Path string
}

func (c FirefoxConnector) Type() string {
	return "firefox"
}

func (c FirefoxConnector) Connect(url string, config *ConnectorConfig) chan bool {
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
		args = append(args, "--new-instance", url)
	}
	args = append(args, "--remote-debugging-port", strconv.Itoa(config.Port))

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

func (c FirefoxConnector) OnDisconnect(cb func(), config *ConnectorConfig) chan bool {
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
