package main

import service "golang.org/x/sys/unix"

func StartPalMonAgent() {
	Info.Println("Starting PalMon Agent")
	service.Exit(8)
}
