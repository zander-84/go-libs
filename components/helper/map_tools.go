package helper

import "sort"

var mapTool = NewMapTool()

type MapTool struct{}

func NewMapTool() *MapTool { return new(MapTool) }

func GetMapTool() *MapTool { return mapTool }

type MapStruct struct {
	Key string
	Val interface{}
}

//MapSort 按key排序
func (t *MapTool) MapSort(m map[string]interface{}) (res []MapStruct) {
	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, key := range keys {
		res = append(res, MapStruct{
			Key: key,
			Val: m[key],
		})
	}
	return res
}
