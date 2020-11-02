package link

type doubleLinkNode struct {
	value interface{}
	prev  *doubleLinkNode
	next  *doubleLinkNode
}

func newDoubleLinkNode(value interface{}) *doubleLinkNode {
	return &doubleLinkNode{value: value, prev: nil, next: nil}
}

func (this *doubleLinkNode) Value() interface{} {
	return this.value
}

func (this *doubleLinkNode) Next() *doubleLinkNode {
	return this.next
}

func (this *doubleLinkNode) Prev() *doubleLinkNode {
	return this.prev
}
