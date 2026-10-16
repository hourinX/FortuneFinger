package pattern

import "regexp"

// Mobile 匹配中国大陆手机号。
var Mobile = regexp.MustCompile(`^1[3-9]\d{9}$`)
