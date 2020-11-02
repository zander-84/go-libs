package link

import (
	"fmt"
	"sync"
)

type doubleLinkList struct {
	head   *doubleLinkNode
	tail   *doubleLinkNode
	length int
	mutex  sync.Mutex
}

func newDoubleLinkList() LinkList {
	head := newDoubleLinkNode(nil)
	return &doubleLinkList{head: head, tail: head, length: 0}
}

func (this *doubleLinkList) Len() int {
	return this.length
}

func (this *doubleLinkList) Peek(index int) interface{} {
	if index > this.length-1 || index < 0 {
		return nil
	}

	if this.length/2 >= index {
		bak := this.head
		for index > -1 {
			bak = bak.next
			index--
		}
		return bak.value
	} else {
		bak := this.tail
		gap := this.length - 1 - index

		for gap > 0 {
			bak = bak.prev
			gap--
		}
		return bak.value
	}
}

func (this *doubleLinkList) Get(index int) interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if index > this.length-1 || index < 0 {
		return nil
	}
	var bak *doubleLinkNode
	if this.length/2 >= index {
		bak = this.head
		gap := index + 1
		for gap > 0 {
			bak = bak.next
			gap--
		}
	} else {
		bak = this.tail
		gap := this.length - 1 - index

		for gap > 0 {
			bak = bak.prev
			gap--
		}
	}
	val := bak.value

	//最后一个
	if bak.next == nil {
		bak.prev.next = nil
		this.tail = bak.prev
	} else {
		bak.next.prev = bak.prev
		bak.prev.next = bak.next
	}
	this.length--
	return val

}

func (this *doubleLinkList) Del(index int) bool {
	d := this.Get(index)
	if d == nil {
		return false
	}
	return true
}

// 头 插入
func (this *doubleLinkList) Unshift(val interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	node := newDoubleLinkNode(val)
	if this.head.next == nil {
		this.head.next = node
		node.prev = this.head
		node.next = nil
		this.tail = node
	} else {
		//交换位置注意顺序
		this.head.next.prev = node
		node.prev = this.head
		node.next = this.head.next
		this.head.next = node
	}
	this.length++

}

// 尾 插入
func (this *doubleLinkList) Push(val interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	node := newDoubleLinkNode(val)

	if this.head.next == nil {
		this.head.next = node
		node.prev = this.head
		node.next = nil
		this.tail = node
	} else {
		bak := this.head
		for bak.next != nil {
			bak = bak.next
		}
		bak.next = node
		node.prev = bak
		this.tail = node
	}
	this.length++
}

// 头 移出
func (this *doubleLinkList) Shift() interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.head.next == nil {
		return nil
	} else {
		val := this.head.next.value

		if this.head.next.next == nil {
			this.head.next = nil
			this.tail = this.head
		} else {
			this.head.next.next.prev = this.head
			this.head.next = this.head.next.next
		}

		this.length--
		return val
	}
}

// 尾 移出
func (this *doubleLinkList) Pop() interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.head.next == nil {
		return nil
	} else {
		value := this.tail.value
		this.tail.prev.next = nil
		this.tail = this.tail.prev
		this.length--
		return value
	}
}

func (this *doubleLinkList) String() string {
	var listString1 string
	var listString2 string
	p := this.head
	for p.next != nil {
		listString1 += fmt.Sprintf("%v-->", p.next.value)
		p = p.next
	}
	listString1 += fmt.Sprintf("nil\n")
	listString1 += fmt.Sprintf("nil")
	for p != this.head {
		listString2 += fmt.Sprintf("<--%v", p.value)
		p = p.prev
	}

	return listString1 + listString2 + "\n"
}
