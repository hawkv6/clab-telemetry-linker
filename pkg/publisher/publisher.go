package publisher

var subsystem = "publisher"

type Publisher interface {
	Init() error
	Start()
	Stop() error
}
