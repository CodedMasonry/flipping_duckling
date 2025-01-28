//go:build windows

package main

import "time"

func main() {
	for {
		PatchAllPowershells("powershell.exe")
		time.Sleep(ptchdrq * time.Millisecond)
	}
}
