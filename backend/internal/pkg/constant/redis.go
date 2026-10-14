package constant

import "time"

const (
	MenuCreateDebounceKey     = "menu:create:debounce"
	MenuUpdateDebounceKey     = "menu:update:debounce"
	MenuDeleteDebounceKey     = "menu:delete:debounce"
	RoleCreateDebounceKey     = "role:create:debounce"
	RoleUpdateDebounceKey     = "role:update:debounce"
	RoleDeleteDebounceKey     = "role:delete:debounce"
	RolePermissionDebounceKey = "role:permission:debounce"
	UserCreateDebounceKey     = "user:create:debounce"
	UserUpdateDebounceKey     = "user:update:debounce"
	UserDeleteDebounceKey     = "user:delete:debounce"
	UserRoleDebounceKey       = "user:role:debounce"
	OrgCreateDebounceKey      = "org:create:debounce"
	OrgUpdateDebounceKey      = "org:update:debounce"
	OrgDeleteDebounceKey      = "org:delete:debounce"
	DictTypeCreateDebounceKey = "dict:type:create:debounce"
	DictTypeUpdateDebounceKey = "dict:type:update:debounce"
	DictTypeDeleteDebounceKey = "dict:type:delete:debounce"
	DictDataCreateDebounceKey = "dict:data:create:debounce"
	DictDataUpdateDebounceKey = "dict:data:update:debounce"
	DictDataDeleteDebounceKey = "dict:data:delete:debounce"
)

const CommonDebounceTTL = 3 * time.Second

const (
	CaptchaKeyPrefix      = "captcha:"
	RefreshTokenKeyPrefix = "auth:refresh:"
)
