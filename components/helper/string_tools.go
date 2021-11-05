package helper

import (
	"strings"
)

var stringTool = NewStringTool()

type StringTool struct{}

func NewStringTool() *StringTool { return new(StringTool) }

func GetStringTool() *StringTool { return stringTool }

// GetFirstString 没匹配到返回所有
func (t *SliceTool) GetFirstString(s string, substr string) string {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return s
	}
	return s[:index]
}

func (t *StringTool) GetFirstInt(s string, substr string) int {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return 0
	}
	return defaultConv.ShouldStoI(s[:index])
}

// GetLastString 没匹配到返回空
func (t *StringTool) GetLastString(s string, substr string) string {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return ""
	}
	return s[index+len(substr):]
}

func (t *StringTool) GetLastInt(s string, substr string) int {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return 0
	}
	return defaultConv.ShouldStoI(s[index+len(substr):])
}

func (t *StringTool) GetLastInt64(s string, substr string) int64 {
	index := strings.LastIndex(s, substr)
	if index < 0 {
		return 0
	}
	return defaultConv.ShouldStoI64(s[index+len(substr):])
}
