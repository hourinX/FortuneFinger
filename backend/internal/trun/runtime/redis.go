package runtime

import (
	"errors"

	"fortunefinger/backend/internal/pkg/redisc"
)

func (e *Environment) Redis() (*redisc.Client, error) {
	if e == nil || e.redis == nil {
		return nil, errors.New("trun: 运行环境尚未初始化")
	}
	return e.redis, nil
}
