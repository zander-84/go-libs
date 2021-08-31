package helper

import (
	"fmt"
	"strconv"
)

var defaultConv = NewConv()

type Conv struct{}

func NewConv() *Conv { return new(Conv) }

func GetConv() *Conv { return defaultConv }

//ShouldStoI 字符串 转 数字
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
	int10, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return int10
}

//ShouldUtoS
//只限10进制
func (c *Conv) ShouldUtoS(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

//ShouldU32toS
//只限10进制
func (c *Conv) ShouldU32toS(u uint32) string {
	return strconv.FormatUint(uint64(u), 10)
}

//ShouldU64toS
//只限10进制
func (c *Conv) ShouldU64toS(u uint64) string {
	return strconv.FormatUint(u, 10)
}

//ShouldStoB
//字符串bool转系统bool
//备注：
//True类型:"1", "t", "T", "true", "TRUE", "True"
//False类型:"0", "f", "F", "false", "FALSE", "False"
func (c *Conv) ShouldStoB(b string) bool {
	bo, err := strconv.ParseBool(b)
	if err != nil {
		return false
	}
	return bo
}

//ShouldStoU
//字符串转无符号整型
//备注：长度不能超过系统类型长度/其他非数字/为空，否则返回0
//10进制
//u 无符号数字字符串
//返回结果uint
func (c *Conv) ShouldStoU(u string) uint {
	uInt64, err := strconv.ParseUint(u, 10, 0)
	if err != nil {
		return 0
	}
	return uint(uInt64)
}

//ShouldStoU32
//字符串转无符号整型
//备注：长度不能超过系统类型长度/其他非数字/为空，否则返回0
//10进制
//u 无符号数字字符串
//返回结果uint32
func (c *Conv) ShouldStoU32(u string) uint32 {
	uInt64, err := strconv.ParseUint(u, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(uInt64)
}

//ShouldStoU64
//字符串转无符号整型
//备注：长度不能超过系统类型长度/其他非数字/为空，否则返回0
//10进制
//u 无符号数字字符串
//返回结果uint64
func (c *Conv) ShouldStoU64(u string) uint64 {
	uInt64, err := strconv.ParseUint(u, 10, 64)
	if err != nil {
		return 0
	}
	return uInt64
}

// 数字 转 字符串
func (c *Conv) ShouldItoS(i int) string {
	return fmt.Sprintf("%d", i)
}

func (c *Conv) ShouldI32toS(i int32) string {
	return fmt.Sprintf("%d", i)
}

func (c *Conv) ShouldI64toS(i int64) string {
	return fmt.Sprintf("%d", i)
}

// 字符串 转 浮点
func (c *Conv) ShouldStoF(s string) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return n
}

// 浮点 转 字符串
func (c *Conv) ShouldFtoS(f float64) string {
	return fmt.Sprintf("%f", f)
}

// 定点浮点 四舍五入
func (c *Conv) ShouldDecimal(f float64, len int) float64 {
	format := fmt.Sprintf("%%.%df", len)
	value, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return value
}
