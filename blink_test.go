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
	view.LoadURL("https://open.163.com/appdownload")
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
	<-exit
}
