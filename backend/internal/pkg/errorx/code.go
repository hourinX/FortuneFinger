package errorx

import "net/http"

type Code int

const (
	InvalidParameter             Code = 1
	Unauthorized                 Code = 2
	Forbidden                    Code = 3
	NotFound                     Code = 4
	Conflict                     Code = 5
	TooManyRequests              Code = 6
	InternalServer               Code = 7
	DatabaseError                Code = 8
	ServiceUnavailable           Code = 9
	EnhancedVerificationRequired Code = 10
	RequestTimeout               Code = 13
)

type definition struct {
	status  int
	message string
}

var definitions = map[Code]definition{
	InvalidParameter: {
		status:  http.StatusBadRequest,
		message: "请求参数有误",
	},
	Unauthorized: {
		status:  http.StatusUnauthorized,
		message: "身份认证失败，请重新登录",
	},
	Forbidden: {
		status:  http.StatusForbidden,
		message: "当前账号无权执行该操作",
	},
	NotFound: {
		status:  http.StatusNotFound,
		message: "请求的资源不存在",
	},
	Conflict: {
		status:  http.StatusConflict,
		message: "当前资源状态发生冲突",
	},
	TooManyRequests: {
		status:  http.StatusTooManyRequests,
		message: "操作过于频繁，请稍后重试",
	},
	InternalServer: {
		status:  http.StatusInternalServerError,
		message: "系统繁忙，请稍后重试",
	},
	DatabaseError: {
		status:  http.StatusInternalServerError,
		message: "数据服务异常，请稍后重试",
	},
	ServiceUnavailable: {
		status:  http.StatusServiceUnavailable,
		message: "服务暂时不可用，请稍后重试",
	},
	EnhancedVerificationRequired: {
		status:  http.StatusPreconditionRequired,
		message: "请先完成增强验证",
	},
	RequestTimeout: {
		status:  http.StatusGatewayTimeout,
		message: "当前服务响应已超时，请稍后重试",
	},
}

func lookup(code Code) (Code, definition) {
	def, exists := definitions[code]
	if exists {
		return code, def
	}
	return InternalServer, definitions[InternalServer]
}

func HTTPStatus(code Code) int {
	_, def := lookup(code)
	return def.status
}
