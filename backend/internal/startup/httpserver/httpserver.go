package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	readHeaderTimeout = 3 * time.Second  // 读取请求头的最长时间
	readTimeout       = 60 * time.Second // 读取完整请求的最长时间
	writeTimeout      = 30 * time.Second // 写入完整响应的最长时间
	idleTimeout       = 60 * time.Second // Keep-Alive空闲连接的最长保留时间
	shutdownTimeout   = 20 * time.Second // 优雅关闭的最长等待时间，必须小于部署环境的终止宽限期
	maxHeaderBytes    = 64 << 10         // 请求头最大字节数，当前为64KiB
)

func New(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}

func Shutdown(ctx context.Context, srv *http.Server) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()
	shutdownErr := srv.Shutdown(shutdownCtx)
	if shutdownErr == nil {
		return nil
	}

	closeErr := srv.Close()
	if closeErr != nil {
		return errors.Join(
			fmt.Errorf("httpserver: 优雅关闭HTTP服务失败: %w", shutdownErr),
			fmt.Errorf("httpserver: 强制关闭HTTP服务失败: %w", closeErr),
		)
	}
	return fmt.Errorf("httpserver: 优雅关闭HTTP服务失败: %w", shutdownErr)
}
