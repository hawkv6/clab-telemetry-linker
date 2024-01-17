package service

var subsystem = "service"

type Service interface {
	Start()
	Stop() error
}
