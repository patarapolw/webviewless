package webviewless

// cSpell: disable

/*
#cgo linux openbsd freebsd pkg-config: gtk+-3.0 webkit2gtk-4.0
#cgo darwin LDFLAGS: -framework CoreGraphics

#if defined(__APPLE__)
#include <CoreGraphics/CGDisplayConfiguration.h>
int display_width() {
	return CGDisplayPixelsWide(CGMainDisplayID());
}
int display_height() {
	return CGDisplayPixelsHigh(CGMainDisplayID());
}
#elif defined(_WIN32)
#include <wtypes.h>
int display_width() {
	RECT desktop;
	const HWND hDesktop = GetDesktopWindow();
	GetWindowRect(hDesktop, &desktop);
	return desktop.right;
}
int display_height() {
	RECT desktop;
	const HWND hDesktop = GetDesktopWindow();
	GetWindowRect(hDesktop, &desktop);
	return desktop.bottom;
}
#else
#include <gtk/gtk.h>
int display_width() {
	return 0;
}
int display_height() {
	return 0;
}
#endif
*/
import "C"

import (
	"runtime"

	"github.com/webview/webview"
)

// cSpell: enable

type WebviewConfig struct {
	Title            string
	Width            int
	Height           int
	IsMaximized      bool
	Debug            bool
	PreferredBrowser string
	Native           bool
}

// LaunchWebview launch default webview, is blocking
// and must be run on first thread in macOS
func LaunchWebview(url string, config WebviewConfig, onDestroy func()) {
	if !config.Native {
		loc := LocateBrowser(config.PreferredBrowser)
		if loc.Path != "" {
			cfg := ConnectorConfig{
				Width:       config.Width,
				Height:      config.Height,
				IsMaximized: config.IsMaximized,
			}

			c := *loc.Connector
			<-c.Connect(url, &cfg)
			<-c.OnDisconnect(onDestroy, &cfg)
			return
		}
	}

	w := webview.New(config.Debug)
	defer w.Destroy()
	defer onDestroy()

	if config.IsMaximized {
		switch runtime.GOOS {
		case "linux":
			// Call `gtk_window_fullscreen`, convert window to `C.GtkWindow` pointer.
			C.gtk_window_fullscreen((*C.GtkWindow)(w.Window()))
		default:
			w.SetSize(int(C.display_width()), int(C.display_height()), webview.HintNone)
		}
	} else {
		if config.Height == 0 {
			config.Height = 600
		}
		if config.Width == 0 {
			config.Width = 800
		}

		w.SetSize(config.Width, config.Height, webview.HintNone)
	}

	if config.Title != "" {
		w.SetTitle(config.Title)
	}

	w.Navigate(url)
	w.Run()
}
