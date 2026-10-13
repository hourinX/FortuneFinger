package operator

import (
	"context"
	"strconv"

	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/requestctx/userid"
)

func UserID(ctx context.Context) (uint64, error) {
	value := userid.Get(ctx)
	if value == "" {
		return 0, errorx.New(errorx.Unauthorized, "未获取到当前用户信息")
	}
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, errorx.Wrap(errorx.Unauthorized, "当前用户信息无效", err)
	}
	return id, nil
}
