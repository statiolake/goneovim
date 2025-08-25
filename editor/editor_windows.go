package editor

import "C"

import (
	"syscall"
	"unsafe"
)

const (
	INPUT_KEYBOARD        = 1
	KEYEVENTF_EXTENDEDKEY = 0x0001
	KEYEVENTF_KEYUP       = 0x0002
	VK_LMENU              = 0xA4
	MAPVK_VK_TO_VSC       = 0
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procSendInput           = user32.NewProc("SendInput")
	procMapVirtualKeyW      = user32.NewProc("MapVirtualKeyW")
	procIsIconic            = user32.NewProc("IsIconic")
	procIsWindowVisible     = user32.NewProc("IsWindowVisible")
)

type KEYBDINPUT struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type INPUT struct {
	inputType uint32
	ki        KEYBDINPUT
	_padding  [8]byte // Padding to match C struct alignment
}

func GetOpeningFilepath(str *C.char) {
}

func setMyApplicationDelegate() {
}

// setForegroundWindow sets the specified window to the foreground
func setForegroundWindow(hwnd uintptr) bool {
	ret, _, _ := procSetForegroundWindow.Call(hwnd)
	return ret != 0
}

// mapVirtualKey maps a virtual key code to a scan code
func mapVirtualKey(vk uint32, mapType uint32) uint32 {
	ret, _, _ := procMapVirtualKeyW.Call(uintptr(vk), uintptr(mapType))
	return uint32(ret)
}

// sendInput simulates keyboard input
func sendInput(inputs []INPUT) uint32 {
	if len(inputs) == 0 {
		return 0
	}
	ret, _, _ := procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(unsafe.Sizeof(INPUT{})),
	)
	return uint32(ret)
}

// activateWindow brings the window to the foreground
func (e *Editor) activateWindow() {
	hwnd := uintptr(e.window.WinId())

	// Simply using window.Raise() or window.ActivateWindow() is not enough on
	// Windows due to focus stealing prevention. The exact reason is not clear,
	// but they just make the window flash in the taskbar without bringing it to
	// the foreground.
	//
	// The following hacky workaround is based on rust-windowing/winit, which
	// is the backend of Neovide.
	// See: https://github.com/rust-windowing/winit/blob/488c036a05d418e13bbdcdc349e2db2f6b7f58e2/winit-win32/src/window.rs#L1596-L1634

	// ALT key hack to overcome Windows focus stealing protection
	altScanCode := mapVirtualKey(VK_LMENU, MAPVK_VK_TO_VSC)

	inputs := []INPUT{
		{
			inputType: INPUT_KEYBOARD,
			ki: KEYBDINPUT{
				wVk:         VK_LMENU,
				wScan:       uint16(altScanCode),
				dwFlags:     KEYEVENTF_EXTENDEDKEY,
				time:        0,
				dwExtraInfo: 0,
			},
		},
		{
			inputType: INPUT_KEYBOARD,
			ki: KEYBDINPUT{
				wVk:         VK_LMENU,
				wScan:       uint16(altScanCode),
				dwFlags:     KEYEVENTF_EXTENDEDKEY | KEYEVENTF_KEYUP,
				time:        0,
				dwExtraInfo: 0,
			},
		},
	}

	// Simulate ALT key press and release
	sendInput(inputs)

	// Set foreground window
	setForegroundWindow(hwnd)
}
