package queue

type Queue interface {
	// 大小
	Size() int

	// 是否空
	IsEmpty() bool

	// 压入
	Push(data interface{})

	// 压入带优先级的
	PushPriority(data interface{}, priority int)

	// 弹出
	Pop() interface{}

	// 弹出
	PopPriority() (data interface{}, priority int)

	// 清除
	Clear()
}

const (
	ArrayQ = iota
	LinkListQ
	MaxHeapQ
	MinHeapQ
	MaxPriorityQ
	MinPriorityQ
)

func NewQueue(typ int) Queue {
	switch typ {

	// 数组队列 性能一般
	case ArrayQ:
		return newArrayQueue()

	// 列表队列 性能好
	case LinkListQ:
		return newLinkListQueue()

	// 最小队列 性能一般
	case MinHeapQ:
		return newMinHeapQueue()

	// 最大队列 性能一般
	case MaxHeapQ:
		return newMaxHeapQueue()

	// 最小优先队列 性能好
	case MinPriorityQ:
		return newMinPriorityQueue()

	// 最大优先队列 性能好
	case MaxPriorityQ:
		return newMaxPriorityQueue()

	default:
		return newLinkListQueue()
	}
}
