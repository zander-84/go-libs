package nsq

import (
	"github.com/nsqio/go-nsq"
	"time"
)

func GetProducer(conf *nsq.Config, addr string) (producer *nsq.Producer, err error) {
	if conf == nil {
		conf = nsq.NewConfig()
	}

	if producer, err = nsq.NewProducer(addr, conf); err != nil {
		return nil, err
	}

	if err := producer.Ping(); err != nil {
		return nil, err
	}

	return producer, nil
}

func StopProducer(producer *nsq.Producer) {
	producer.Stop()
}

func Publish(producer *nsq.Producer, topic string, body []byte) error {
	return producer.Publish(topic, body)
}

func MultiPublish(producer *nsq.Producer, topic string, body [][]byte) error {
	return producer.MultiPublish(topic, body)
}

//推动到延迟队列  适用于很多定时场景
func DeferredPublish(producer *nsq.Producer, topic string, delay time.Duration, body []byte) error {
	return producer.DeferredPublish(topic, delay, body)
}
