package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

var tool = NewTool()

type Tool struct{}

func NewTool() *Tool { return new(Tool) }

func GetTool() *Tool { return tool }

func (t *Tool) InStringArr(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}
func (t *Tool) InIntArr(s int, ss []int) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

func (t *Tool) GetLastIndexString(s string, substr string) string {
	return s[strings.LastIndex(s, substr)+len(substr):]
}

func (t *Tool) GetLastIndexInt(s string, substr string) int {
	return defaultConv.ShouldStoI(s[strings.LastIndex(s, substr)+len(substr):])
}

func (t *Tool) GetLastIndexInt64(s string, substr string) int64 {
	return defaultConv.ShouldStoI64(s[strings.LastIndex(s, substr)+len(substr):])
}

func PrettyPrint(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(v)
		return
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Println(v)
		return
	}

	fmt.Println(out.String())
}
