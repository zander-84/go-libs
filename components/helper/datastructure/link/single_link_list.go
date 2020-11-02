package link

import (
	"fmt"
	"sync"
)

// 单链适用场景 栈，队列场景选择双链
type singleLinkList struct {
	head   *singleLinkNode
	length int
	mutex  sync.Mutex
}

func newSingleLinkList() LinkList {
	head := newSingleLinkNode(nil)
	return &singleLinkList{head: head, length: 0}
}

func (this *singleLinkList) Len() int {
	return this.length
}

// 头 插入
func (this *singleLinkList) Unshift(val interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	node := newSingleLinkNode(val)
	if this.head.next == nil {
		this.head.next = node
		node.next = nil
	} else {
		node.next = this.head.next
		this.head.next = node
	}
	this.length++

}

// 尾 插入
func (this *singleLinkList) Push(val interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	node := newSingleLinkNode(val)

	if this.head.next == nil {
		this.head.next = node
		node.next = nil
	} else {
		bak := this.head
		for bak.next != nil {
			bak = bak.next
		}
		bak.next = node
	}
	this.length++
}

// 头 移出
func (this *singleLinkList) Shift() interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.head.next == nil {
		return nil
	} else {
		val := this.head.next.value
		this.head.next = this.head.next.next
		this.length--
		return val
	}
}

// 尾 移出
func (this *singleLinkList) Pop() interface{} {
	return this.Get(this.length - 1)
}
func (this *singleLinkList) Peek(index int) interface{} {
	if index > this.length-1 || index < 0 {
		return nil
	} else {
		phead := this.head
		for index > -1 {
			phead = phead.next
			index--
		}
		return phead.value
	}
}
func (this *singleLinkList) DeleteNode(dest interface{}) bool {
	if dest == nil {
		return false
	} else {
		phead := this.head
		for phead.next != nil && phead.next != dest {
			phead = phead.next
		}
		if phead.next.value == dest {
			phead.next = phead.next.next
			this.length--
			return true
		} else {
			return false
		}
	}
}

func (this *singleLinkList) Get(index int) interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if index > this.length-1 || index < 0 {
		return nil
	} else {
		phead := this.head
		for index > 0 {
			phead = phead.next
			index--
		}
		val := phead.next.value
		phead.next = phead.next.next
		this.length--
		return val
	}

}
func (this *singleLinkList) Del(index int) bool {
	d := this.Get(index)
	if d == nil {
		return false
	}
	return true
}

// 需要修改
func (this *singleLinkList) InserBefore(dest interface{}, node *singleLinkNode) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	phead := this.head
	isFind := false
	for phead.next != nil {
		if phead.next.value == dest {
			isFind = true
			break
		}
		phead = phead.next
	}
	if isFind {
		node.next = phead.next
		phead.next = node
		this.length++
		return true
	} else {
		return false
	}
}

// 需要修改
func (this *singleLinkList) InserAfter(dest interface{}, node *singleLinkNode) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	phead := this.head
	isFind := false
	for phead.next != nil {
		if phead.next.value == dest {
			isFind = true
			break
		}
		phead = phead.next
	}
	if isFind {
		node.next = phead.next.next
		phead.next.next = node
		this.length++
		return true
	} else {
		return false
	}
}

func (this *singleLinkList) String() string {
	var listString string
	p := this.head
	for p.next != nil {
		listString += fmt.Sprintf("%v-->", p.next.value)
		p = p.next
	}
	listString += fmt.Sprintf("nil")
	return listString
}
