package publisher

import "github.com/sirupsen/logrus"

var subsystem = "publisher"

type Publisher interface {
	Start() error
}

type DefaultPublisher struct {
	log *logrus.Entry
}

func NewDefaultPublisher() *DefaultPublisher {
	return &DefaultPublisher{
		log: logrus.WithField("subsystem", subsystem),
	}
}

func (publisher *DefaultPublisher) Start() error {
	publisher.log.Infoln("Start Publisher")
	return nil
}
