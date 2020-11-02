package queue

import (
	"github.com/zander-84/go-libs/components/helper/datastructure/link"
	"sync"
)

type linkListQueue struct {
	data  link.LinkList
	size  int
	mutex sync.Mutex
}

func newLinkListQueue() Queue {
	this := new(linkListQueue)
	this.Clear()
	return this
}

func (this *linkListQueue) Size() int {
	return this.size
}

func (this *linkListQueue) IsEmpty() bool {
	return this.size == 0
}
func (this *linkListQueue) Push(data interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.data.Unshift(data)
	this.size++
}

func (this *linkListQueue) PushPriority(data interface{}, priority int) {
	this.Push(data)
}

func (this *linkListQueue) Pop() interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.IsEmpty() {
		return nil
	} else {
		this.size--
		return this.data.Pop()
	}
}

func (this *linkListQueue) PopPriority() (data interface{}, priority int) {
	return this.Pop(), 0
}

func (this *linkListQueue) Clear() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.data = link.NewLinkList(link.DoubleLink)
	this.size = 0
}
