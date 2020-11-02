package queue

import "sync"

type arrayQueue struct {
	data  []interface{}
	size  int
	mutex sync.Mutex
}

func newArrayQueue() Queue {
	this := new(arrayQueue)
	this.Clear()
	return this
}

func (this *arrayQueue) Size() int {
	return this.size
}
func (this *arrayQueue) IsEmpty() bool {
	return this.size == 0
}
func (this *arrayQueue) Push(data interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.data = append(this.data, data)
	this.size++
}

func (this *arrayQueue) PushPriority(data interface{}, priority int) {
	this.Push(data)
}

func (this *arrayQueue) Pop() interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.IsEmpty() {
		return nil
	} else {
		data := this.data[0]
		this.data = this.data[1:]
		this.size--
		return data
	}
}

func (this *arrayQueue) PopPriority() (data interface{}, priority int) {
	return this.Pop(), 0
}

func (this *arrayQueue) Clear() {
	this.data = make([]interface{}, 0)
	this.size = 0
}
