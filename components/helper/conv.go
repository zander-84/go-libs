package helper

import (
	"fmt"
	"strconv"
)

type Conv struct{}

func NewConv() *Conv { return new(Conv) }

// 字符串 转 数字
func (c *Conv) ShouldStoI(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

func (c *Conv) ShouldStoI32(s string) int32 {
	int10, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int32(int10)
}

func (c *Conv) ShouldStoI64(s string) int64 {
	int10, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int10
}

// 数字 转 字符串
func (c *Conv) ShouldItos(i int) string {
	return fmt.Sprintf("%d", i)
}

func (c *Conv) ShouldI32tos(i int32) string {
	return fmt.Sprintf("%d", i)
}

func (c *Conv) ShouldI64tos(i int64) string {
	return fmt.Sprintf("%d", i)
}

// 字符串 转 浮点
func (c *Conv) ShouldStof(s string) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return n
}

// 浮点 转 字符串
func (c *Conv) ShouldFtos(f float64) string {
	return fmt.Sprintf("%f", f)
}

// 定点浮点 四舍五入
func (c *Conv) ShouldDecimal(f float64, len int) float64 {
	format := fmt.Sprintf("%%.%df", len)
	value, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return value
}
