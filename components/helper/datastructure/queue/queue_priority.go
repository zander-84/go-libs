package queue

import (
	"sync"
)

// 一级队列新能一般，二级队列性能好
// 适用： 一级队列少数，二级队列多数的场景
type priorityQueue struct {
	size  int
	mutex sync.Mutex
	data  []*priorityItem
	min   bool
}

// 最小 队列
func newMinPriorityQueue() Queue {
	this := new(priorityQueue)
	this.min = false
	this.Clear()
	return this
}

// 最大 队列
func newMaxPriorityQueue() Queue {
	this := new(priorityQueue)
	this.Clear()
	this.min = true
	return this
}

func (this *priorityQueue) Clear() {
	this.data = make([]*priorityItem, 0)
}

func (this *priorityQueue) IsEmpty() bool {
	return this.size == 0
}

func (this *priorityQueue) Size() int {
	return this.size
}

func (this *priorityQueue) Push(it interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	i, ok := it.(*tempPriorityItem)
	if !ok {
		i = newTempPriorityItem(it, 0)
	}

	if this.size == 0 {
		this.data = append(this.data, newPriorityItem(i.value, i.priority))
	} else {
		exist := false
		for _, val := range this.data {
			if val.priority == i.priority {
				exist = true
				val.value.Unshift(i.value)
			}
		}
		if !exist {
			this.data = append(this.data, newPriorityItem(i.value, i.priority))
		}
	}

	this.size++
}

func (this *priorityQueue) PushPriority(data interface{}, priority int) {
	this.Push(newTempPriorityItem(data, priority))
}

func (this *priorityQueue) Less(a, b *priorityItem) bool {
	if this.min {
		return a.Less(b)
	} else {
		return b.Less(a)
	}
}

func (this *priorityQueue) Pop() interface{} {
	data, _ := this.PopPriority()
	return data
}

func (this *priorityQueue) PopPriority() (data interface{}, priority int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.Size() == 0 {
		return nil, 0
	}

	this.heapSort(this.data)
	data = this.data[0].value.Pop()
	priority = this.data[0].priority
	if this.data[0].value.Len() == 0 {
		this.data = this.data[1:]
	}

	this.size--
	return data, priority
}

// 弹出极值
func (this *priorityQueue) heapSort(arr []*priorityItem) {
	length := len(this.data)
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
