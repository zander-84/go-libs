package queue

import (
	"sync"
)

// 堆队列  不适合大数据场景
type heapQueue struct {
	lock sync.Mutex
	data []heapItem
	min  bool
}

// 最小堆
func newMinHeapQueue() Queue {
	this := new(heapQueue)
	this.min = false
	this.Clear()
	return this
}

// 最大堆
func newMaxHeapQueue() Queue {
	this := new(heapQueue)
	this.Clear()
	this.min = true
	return this
}

func (this *heapQueue) Clear() {
	this.data = make([]heapItem, 0)
}

func (this *heapQueue) IsEmpty() bool {
	return len(this.data) == 0
}

func (this *heapQueue) Size() int {
	return len(this.data)
}

func (this *heapQueue) Push(it interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	i, ok := it.(heapItem)
	if ok {
		this.data = append(this.data, i)
	}
}

func (this *heapQueue) PushPriority(data interface{}, priority int) {
	this.Push(newHeapItem(data, priority))
}

func (this *heapQueue) Less(a, b heapItem) bool {
	if this.min {
		return a.Less(b)
	} else {
		return b.Less(a)
	}
}

func (this *heapQueue) Pop() interface{} {
	data, _ := this.PopPriority()
	return data
}

func (this *heapQueue) PopPriority() (data interface{}, priority int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.Size() == 0 {
		return nil, 0
	}

	this.heapSort(this.data)
	data = this.data[0].value
	priority = this.data[0].priority
	this.data = this.data[1:]

	return data, priority
}

// 弹出极值
func (this *heapQueue) heapSort(arr []heapItem) {
	length := this.Size()
	if length <= 1 {
		return
	}

	depth := length/2 - 1
	for i := depth; i >= 0; i-- {
		var tmp int
		topmax := i
		leftchild := 2*i + 1
		rightchild := 2*i + 2

		if rightchild <= length-1 && this.Less(arr[rightchild], arr[leftchild]) {
			tmp = rightchild
		} else {
			tmp = leftchild
		}

		if this.Less(arr[tmp], arr[topmax]) {
			topmax = tmp
		}

		if topmax != i {
			arr[i], arr[topmax] = arr[topmax], arr[i]
		}
	}

}
