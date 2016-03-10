package service_control

import (
	"golang.org/x/sys/windows/svc"
)

type ServiceControlWindows struct {

}

func (sc *ServiceControlWindows) Install(svcName, svcDescription string) error {
	return installService(svcName, svcDescription)
}

func (sc *ServiceControlWindows) Remove(svcName string) error {
	return removeService(svcName)
}

func (sc *ServiceControlWindows) Start(svcName string) error {
	return startService(svcName)
}

func (sc *ServiceControlWindows) Stop(svcName string) error {
	return controlService(svcName, svc.Stop, svc.Stopped)
}