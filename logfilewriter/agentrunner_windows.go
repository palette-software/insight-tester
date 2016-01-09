package main

import service "golang.org/x/sys/windows/svc"

type PalMonHandler struct {

}

func (pmh PalMonHandler) Execute(args []string, r <-chan ChangeRequest, s chan<- Status) (svcSpecificEC bool, exitCode uint32) {

}

func StartPalMonAgent() {
	Info.Println("Starting PalMon Agent")
	var handler = PalMonHandler{}
	service.Run("PalMonAgent", handler)
}

