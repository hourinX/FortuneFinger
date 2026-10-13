package mask

import "strings"

// Name 隐藏姓名中间部分，两个字符保留首字符，三个及以上字符保留首尾字符。
func Name(value string) string {
	chars := []rune(value)
	switch len(chars) {
	case 0:
		return ""
	case 1:
		return "*"
	case 2:
		return string(chars[:1]) + "*"
	default:
		return string(chars[:1]) + "*" + string(chars[len(chars)-1:])
	}
}

// Mobile 隐藏手机号中间部分，保留前三位和后四位。
func Mobile(value string) string {
	chars := []rune(value)
	length := len(chars)
	if length == 0 {
		return ""
	}
	if length <= 2 {
		return strings.Repeat("*", length)
	}
	if length <= 7 {
		return string(chars[:1]) + strings.Repeat("*", length-2) + string(chars[length-1:])
	}
	return string(chars[:3]) + strings.Repeat("*", length-7) + string(chars[length-4:])
}
