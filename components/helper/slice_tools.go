package helper

var sliceTool = NewSliceTool()

type SliceTool struct{}

func NewSliceTool() *SliceTool { return new(SliceTool) }

func GetSliceTool() *SliceTool { return sliceTool }

func (t *SliceTool) InStringArr(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}
func (t *SliceTool) InIntArr(s int, ss []int) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

func (t *SliceTool) InInt32Arr(s int32, ss []int32) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

func (t *SliceTool) InInt64Arr(s int64, ss []int64) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

//ShouldI64SliceToStr []int64 转 []string
func (t *SliceTool) ShouldI64SliceToStr(i []int64) []string {
	s := make([]string, 0)
	for _, v := range i {
		s = append(s, GetConv().ShouldI64toS(v))
	}
	return s
}

//ShouldStrSliceToI64 []string 转 []int64
func (t *SliceTool) ShouldStrSliceToI64(s []string) []int64 {
	i := make([]int64, 0)
	for _, v := range s {
		i = append(i, GetConv().ShouldStoI64(v))
	}
	return i
}

//ShouldI32SliceToStr []int32 转 []string
func (t *SliceTool) ShouldI32SliceToStr(i []int32) []string {
	s := make([]string, 0)
	for _, v := range i {
		s = append(s, GetConv().ShouldI32toS(v))
	}
	return s
}

//ShouldStrSliceToI32 []string 转 []int32
func (t *SliceTool) ShouldStrSliceToI32(s []string) []int32 {
	i := make([]int32, 0)
	for _, v := range s {
		i = append(i, GetConv().ShouldStoI32(v))
	}
	return i
}

//ShouldISliceToStr []int 转 []string
func (t *SliceTool) ShouldISliceToStr(i []int) []string {
	s := make([]string, 0)
	for _, v := range i {
		s = append(s, GetConv().ShouldItoS(v))
	}
	return s
}

//ShouldStrSliceToI []string 转 []int
func (t *SliceTool) ShouldStrSliceToI(s []string) []int {
	i := make([]int, 0)
	for _, v := range s {
		i = append(i, GetConv().ShouldStoI(v))
	}
	return i
}
