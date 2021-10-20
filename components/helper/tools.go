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

// GetLastIndexString 没匹配到返回空
func (t *Tool) GetLastIndexString(s string, substr string) string {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return ""
	}
	return s[index+len(substr):]
}

// GetIndexString 没匹配到返回所有
func (t *Tool) GetIndexString(s string, substr string) string {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return s
	}
	return s[:index]
}

func (t *Tool) GetIndexInt(s string, substr string) int {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return 0
	}
	return defaultConv.ShouldStoI(s[:index])
}

func (t *Tool) GetLastIndexInt(s string, substr string) int {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return 0
	}
	return defaultConv.ShouldStoI(s[index+len(substr):])
}

func (t *Tool) GetLastIndexInt64(s string, substr string) int64 {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return 0
	}
	return defaultConv.ShouldStoI64(s[index+len(substr):])
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
