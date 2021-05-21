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
func LocateChrome() BrowserProtocol {
	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/usr/bin/ungoogled-chromium",
			"/usr/bin/ungoogled-chromium-browser",
		}
	case "windows":
		paths = []string{
			os.Getenv("LocalAppData") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("ProgramFiles") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("ProgramFiles(x86)") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("LocalAppData") + "/Chromium/Application/chrome.exe",
			os.Getenv("ProgramFiles") + "/Chromium/Application/chrome.exe",
			os.Getenv("ProgramFiles(x86)") + "/Chromium/Application/chrome.exe",
			// TODO: Add Chromium
		}
	default:
		paths = []string{
			"/usr/bin/google-chrome-stable",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
			"/usr/bin/ungoogled-chromium",
			"/usr/bin/ungoogled-chromium-browser",
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
			Type:      "chrome",
			Connector: &c,
		}
	}

	return BrowserProtocol{}
}

// ChromeConnector connector for Google Chrome
type ChromeConnector struct {
	Path string
}

func (c ChromeConnector) Type() string {
	return "chrome"
}

func (c ChromeConnector) Connect(url string, config *ConnectorConfig) chan bool {
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

func (ChromeConnector) OnDisconnect(cb func(), config *ConnectorConfig) chan bool {
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
