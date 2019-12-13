package blink

import (
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
	view.LoadURL("http://127.0.0.1")
	view.SetWindowTitle("test")
	view.MoveToCenter()
	view.ShowWindow()
	view.On("destroy", func(_ *WebView) {
		close(exit)
	})
	<-exit
}
