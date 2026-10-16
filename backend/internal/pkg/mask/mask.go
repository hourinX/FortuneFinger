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

// Email 隐藏邮箱名称中间部分，保留域名用于识别邮箱来源。
func Email(value string) string {
	separator := strings.LastIndex(value, "@")
	if separator <= 0 {
		return Name(value)
	}
	name := []rune(value[:separator])
	domain := value[separator:]
	switch len(name) {
	case 0:
		return "*" + domain
	case 1:
		return "*" + domain
	case 2:
		return string(name[:1]) + "*" + domain
	default:
		return string(name[:1]) + strings.Repeat("*", len(name)-2) + string(name[len(name)-1:]) + domain
	}
}
