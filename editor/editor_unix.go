//go:build linux || freebsd
// +build linux freebsd

package editor

import "C"

func GetOpeningFilepath(str *C.char) {
}

func setMyApplicationDelegate() {
}

// activateWindow brings the window to the foreground on Unix systems
func (e *Editor) activateWindow() {
	e.window.Raise()
	e.window.ActivateWindow()
}
