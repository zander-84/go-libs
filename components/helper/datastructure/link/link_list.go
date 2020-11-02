package link

type LinkList interface {
	Len() int

	// 查看  不适合大数据
	Peek(index int) interface{}

	// 获取  不适合大数据
	Get(index int) interface{}

	Del(index int) bool

	// 头部 插入
	Unshift(value interface{})

	// 尾部 插入
	Push(value interface{})

	// 头部 移出
	Shift() interface{}

	// 尾部 移出
	Pop() interface{}

	// 打印
	String() string
}

const (
	SingleLink = iota
	DoubleLink
)

func NewLinkList(which int) LinkList {
	if which == SingleLink {
		return newSingleLinkList()
	}
	return newDoubleLinkList()
}
