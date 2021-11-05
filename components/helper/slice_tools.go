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
