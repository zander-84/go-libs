package queue

import "github.com/zander-84/go-libs/components/helper/datastructure/link"

type tempPriorityItem struct {
	value    interface{}
	priority int
}

func newTempPriorityItem(value interface{}, priority int) *tempPriorityItem {
	return &tempPriorityItem{
		value:    value,
		priority: priority,
	}
}

type priorityItem struct {
	value    link.LinkList
	priority int
}

func newPriorityItem(value interface{}, priority int) *priorityItem {
	lk := link.NewLinkList(link.DoubleLink)
	lk.Push(value)
	return &priorityItem{
		value:    lk,
		priority: priority,
	}
}
func (this *priorityItem) Less(than *priorityItem) bool {
	return this.priority > than.priority
}
