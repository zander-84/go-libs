package message

type Message interface {

	// 队列长度
	Size() int

	// 没有优先级生产
	Produce(data interface{})

	// 带有限级别生产
	ProducePriority(data interface{}, priority int)

	// 消费
	Consume(func(message *Messages) error)

	// 消费失败队列
	ConsumeFailQueue(func(message *Messages) error)
}

type Messages struct {
	Data     interface{} // 数据
	Priority int         // 优先级
	isLock   bool        // 锁  不在clean范围内
	consumed bool        // 是否消费
}

func newMessages() *Messages {
	return &Messages{}
}

func (this *Messages) clean() {
	this.Priority = 0
	this.Data = nil
	this.consumed = false
}

func (this *Messages) Finish() {
	this.consumed = true
}
