package blink

import (
	"fmt"
	"log"
	"testing"
)

func TestBlink(t *testing.T) {
	exit := make(chan bool)

	SetDebugMode(true)
	errv := InitBlink()
	if errv != nil {
		log.Fatal(errv)
	}
	view := NewWebView(false, 1366, 920)

	view.SetDebugConfig("wakeMinInterval", "1")
	view.SetDebugConfig("drawMinInterval", "1")
	view.SetDebugConfig("antiAlias", "1")
	view.SetUserAgent("Blink Test")

	view.LoadURL("http://127.0.0.1:8888/agent")
	view.SetWindowTitle("test")
	view.MoveToCenter()
	view.ShowWindow()

	view.On("download", func(v *WebView, url string) {
		fmt.Printf("Cookies: %s\n", view.GetCookie())
		fmt.Printf("URL: %s\n", url)
	})

	view.On("destroy", func(_ *WebView) {
		close(exit)
	})

	fmt.Printf("UserAgent:%s\n", view.GetUserAgent())

	<-exit
}
