package runtime

import (
	"context"
	"errors"
)

func Run(ctx context.Context, configFile string, fn func(*Environment) error) (err error) {
	if fn == nil {
		return errors.New("trun: 运行函数不能为空")
	}

	env, err := New(ctx, configFile)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, env.Close())
	}()

	return fn(env)
}
