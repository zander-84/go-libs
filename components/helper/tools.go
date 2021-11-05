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

// GetFirstString 没匹配到返回所有
func (t *Tool) GetFirstString(s string, substr string) string {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return s
	}
	return s[:index]
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
