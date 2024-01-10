package consumer

var subsystem = "consumer"

type Consumer interface {
	Init() error
	Start()
	Stop() error
}
