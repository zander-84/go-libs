package nsq

import "github.com/nsqio/go-nsq"

func GetConsumer(conf *nsq.Config, topic string, channel string) (*nsq.Consumer, error) {
	if conf == nil {
		conf = nsq.NewConfig()
	}

	return nsq.NewConsumer(topic, channel, conf)
}
func StopConsumer(consumer *nsq.Consumer) {
	consumer.Stop()
}

//  ConnectToNSQLookupds  循环 ConnectToNSQLookupd
func ConsumeByLookupds(consumer *nsq.Consumer, address []string, handler nsq.Handler, concurrency int) error {
	consumer.AddConcurrentHandlers(handler, concurrency)
	return consumer.ConnectToNSQLookupds(address)
}

// ConnectToNSQDs 循环 ConnectToNSQD
func ConsumeByNSQDS(consumer *nsq.Consumer, address []string, handler nsq.Handler, concurrency int) error {
	consumer.AddConcurrentHandlers(handler, concurrency)
	return consumer.ConnectToNSQDs(address)
}

//
func ConsumeByNSQD(consumer *nsq.Consumer, address string, handler nsq.Handler, concurrency int) error {
	consumer.AddConcurrentHandlers(handler, concurrency)
	return consumer.ConnectToNSQD(address)
}

func MultiHandles(consumer *nsq.Consumer, handler nsq.Handler, concurrency int) {
	consumer.AddConcurrentHandlers(handler, concurrency)
}

func Stats(consumer *nsq.Consumer) *nsq.ConsumerStats {
	return consumer.Stats()
}
