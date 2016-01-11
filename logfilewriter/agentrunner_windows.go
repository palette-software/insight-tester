package main

import service "golang.org/x/sys/windows/svc"

type PalMonHandler struct {

}

func (pmh PalMonHandler) Execute(args []string, r <-chan service.ChangeRequest, s chan<- service.Status) (svcSpecificEC bool, exitCode uint32) {
	return true, 0
}

func StartPalMonAgent() {
	Info.Println("Starting PalMon Agent")
	var handler = PalMonHandler{}
	service.Run("PalMonAgent", handler)
}

