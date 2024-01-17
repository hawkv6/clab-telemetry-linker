package processor

var subsystem = "processor"

type Processor interface {
	Start()
	Stop()
}
