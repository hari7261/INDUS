//go:build windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	user32             = syscall.NewLazyDLL("user32.dll")
	setConsoleTitleW   = kernel32.NewProc("SetConsoleTitleW")
	getConsoleMode     = kernel32.NewProc("GetConsoleMode")
	setConsoleMode     = kernel32.NewProc("SetConsoleMode")
	setConsoleCP       = kernel32.NewProc("SetConsoleCP")
	setConsoleOutputCP = kernel32.NewProc("SetConsoleOutputCP")
	allocConsole       = kernel32.NewProc("AllocConsole")
	getConsoleWindow   = kernel32.NewProc("GetConsoleWindow")
	attachConsole      = kernel32.NewProc("AttachConsole")
	freeConsole        = kernel32.NewProc("FreeConsole")
	showWindow         = user32.NewProc("ShowWindow")
)

const (
	ATTACH_PARENT_PROCESS = ^uint32(0) // -1
	SW_SHOW               = 5
	SW_MAXIMIZE           = 3
	utf8CodePage          = 65001
)

func enableConsoleFeatures() {
	// Check if we have a console window
	consoleWindow, _, _ := getConsoleWindow.Call()

	// If no console exists, create a new independent console window
	// This makes INDUS open in its own window when double-clicked
	if consoleWindow == 0 {
		// Try to attach to parent process first (for command-line usage)
		r1, _, _ := attachConsole.Call(uintptr(ATTACH_PARENT_PROCESS))

		// If attach failed (no parent console), create our own
		if r1 == 0 {
			allocConsole.Call()

			// Reopen stdin, stdout, stderr to the new console
			// This is critical for GUI subsystem apps
			stdin, _ := os.Open("CONIN$")
			os.Stdin = stdin

			stdout, _ := os.OpenFile("CONOUT$", os.O_WRONLY, 0)
			os.Stdout = stdout
			os.Stderr = stdout
		}
	}

	// Force UTF-8 so Unicode banner glyphs render correctly in cmd.exe and
	// in the dedicated console allocated by the GUI build.
	setConsoleCP.Call(uintptr(utf8CodePage))
	setConsoleOutputCP.Call(uintptr(utf8CodePage))

	// Enable ANSI escape code support (virtual terminal processing)
	var mode uint32
	handle := syscall.Handle(os.Stdout.Fd())
	getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	// Enable ENABLE_VIRTUAL_TERMINAL_PROCESSING (0x0004)
	// Enable ENABLE_PROCESSED_OUTPUT (0x0001)
	// Enable ENABLE_WRAP_AT_EOL_OUTPUT (0x0002)
	mode |= 0x0004 | 0x0001 | 0x0002
	setConsoleMode.Call(uintptr(handle), uintptr(mode))

	// Also enable for stderr
	handleErr := syscall.Handle(os.Stderr.Fd())
	getConsoleMode.Call(uintptr(handleErr), uintptr(unsafe.Pointer(&mode)))
	mode |= 0x0004 | 0x0001 | 0x0002
	setConsoleMode.Call(uintptr(handleErr), uintptr(mode))
}

func setConsoleTitle(title string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	setConsoleTitleW.Call(uintptr(unsafe.Pointer(titlePtr)))
}
