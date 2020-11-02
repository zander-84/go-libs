package link

type singleLinkNode struct {
	value interface{}
	next  *singleLinkNode
}

func newSingleLinkNode(data interface{}) *singleLinkNode {
	return &singleLinkNode{
		value: data,
		next:  nil,
	}
}

func (this *singleLinkNode) Value() interface{} {
	return this.value
}

func (this *singleLinkNode) Next() *singleLinkNode {
	return this.next
}
