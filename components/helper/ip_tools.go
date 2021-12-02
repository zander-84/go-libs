package helper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ipTool = NewIpTool()

type IpTool struct{}

func NewIpTool() *IpTool { return new(IpTool) }

func GetIpTool() *IpTool { return ipTool }

/**
 * 将IP地址（IPV4）字符串转换为 int类型的数字
 *
 * 思路:将 IP地址（IPV4）的每一段数字转为 8 位二进制数，并将它们放在结果的适当位置上
 *
 * @param IP地址（IPV4） 字符串，如 127.0.0.1
 * @return IP地址（IPV4）  字符串对应的 int值
 */
func (t *IpTool) ipv4ToUint32(ipv4String string) (uint32, error) {
	list := strings.Split(ipv4String, ".")
	if len(list) != 4 {
		return 0, errors.New("ip格式错误")
	}
	result := uint32(0)
	for i, v := range list {
		temp, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return 0, err
		}
		result += uint32(temp) << uint8(24-i*8)
	}
	return result, nil
}

func (t *IpTool) Uint32ToIp(value uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(value>>24), byte(value>>16), byte(value>>8), byte(value))
}
