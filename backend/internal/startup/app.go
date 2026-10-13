package startup

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"fortunefinger/backend/internal/config"
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/startup/dependency"
	"fortunefinger/backend/internal/startup/httpserver"
	"fortunefinger/backend/internal/startup/router"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

type Application struct {
	svcCtx *svc.ServiceContext
	server *http.Server
}

func NewApplication(ctx context.Context, configFile string) (*Application, error) {
	conf, err := config.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("startup: 加载配置失败: %w", err)
	}
	if err := validateService(conf.Service); err != nil {
		return nil, fmt.Errorf("startup: 校验服务配置失败: %w", err)
	}

	svcCtx, err := dependency.Init(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("startup: 初始化服务依赖失败: %w", err)
	}

	addr := net.JoinHostPort(conf.Service.Host, strconv.Itoa(conf.Service.Port))
	handler := mw.Timeout(svcCtx.Logger, router.New(svcCtx))
	return &Application{
		svcCtx: svcCtx,
		server: httpserver.New(addr, handler),
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.server.Addr)
	if err != nil {
		return fmt.Errorf("startup: 监听「%s」失败: %w", a.server.Addr, err)
	}
	a.svcCtx.Logger.Info("HTTP服务启动成功", zap.String("addr", a.server.Addr))

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.server.Serve(listener)
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("startup: HTTP服务运行失败: %w", err)
	case <-ctx.Done():
		if err := httpserver.Shutdown(context.Background(), a.server); err != nil {
			return fmt.Errorf("startup: HTTP服务退出失败: %w", err)
		}
		err := <-errCh
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("startup: HTTP服务运行失败: %w", err)
	}
}

func (a *Application) Close() error {
	if a == nil || a.svcCtx == nil {
		return nil
	}
	return dependency.Close(a.svcCtx)
}

func (a *Application) Logger() *zap.Logger {
	return a.svcCtx.Logger
}

func validateService(conf config.Service) error {
	if strings.TrimSpace(conf.Name) == "" {
		return errors.New("Service.Name不能为空")
	}
	if strings.TrimSpace(conf.Host) == "" {
		return errors.New("Service.Host不能为空")
	}
	if conf.Port <= 0 || conf.Port > 65535 {
		return fmt.Errorf("Service.Port必须在1到65535之间，当前值为「%d」", conf.Port)
	}
	return nil
}
