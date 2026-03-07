//go:build windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	setConsoleTitleW = kernel32.NewProc("SetConsoleTitleW")
	getConsoleMode   = kernel32.NewProc("GetConsoleMode")
	setConsoleMode   = kernel32.NewProc("SetConsoleMode")
)

func enableConsoleFeatures() {
	var mode uint32
	handle := syscall.Handle(os.Stdout.Fd())
	getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	mode |= 0x0004
	setConsoleMode.Call(uintptr(handle), uintptr(mode))
}

func setConsoleTitle(title string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	setConsoleTitleW.Call(uintptr(unsafe.Pointer(titlePtr)))
}
