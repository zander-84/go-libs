package queue

type heapItem struct {
	value    interface{}
	priority int
}

func newHeapItem(value interface{}, priority int) heapItem {
	return heapItem{value: value, priority: priority}
}
func (this heapItem) Less(than heapItem) bool {
	return this.priority > than.priority
}
