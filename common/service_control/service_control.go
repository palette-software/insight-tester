package service_control

type ServiceControl interface {
	Install(svcName, svcDescription string) error
	Remove(svcName string) error
	Start(svcName string) error
	Stop(svcName string) error
}
