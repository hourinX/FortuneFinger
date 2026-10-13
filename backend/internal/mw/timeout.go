package mw

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"

	"go.uber.org/zap"
)

const requestTimeout = 5 * time.Second

func Timeout(log *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		done := make(chan struct{})
		panicCh := make(chan panicResult, 1)
		writer := &timeoutWriter{
			destination: w,
			header:      make(http.Header),
		}

		go func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					panicCh <- panicResult{
						value: recovered,
						stack: debug.Stack(),
					}
				}
			}()
			next.ServeHTTP(writer, r.WithContext(ctx))
			close(done)
		}()

		select {
		case result := <-panicCh:
			panic(result.value)
		case <-done:
			writer.writeResponse()
		case <-ctx.Done():
			writer.cancel(ctx.Err())
			if ctx.Err() != context.DeadlineExceeded {
				// 父Context已取消，客户端不再等待响应，无需继续观察Handler退出。
				return
			}
			go observeHandlerExit(log, r, done, panicCh)

			log.Warn(
				"HTTP请求处理超时",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("timeout", requestTimeout),
			)
			timeoutErr := errorx.New(errorx.RequestTimeout, "")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(errorx.HTTPStatus(timeoutErr.Code()))
			_ = writeTimeoutBody(w, timeoutErr)
		}
	})
}

type panicResult struct {
	value any
	stack []byte
}

type timeoutWriter struct {
	destination http.ResponseWriter
	header      http.Header
	buffer      bytes.Buffer
	status      int
	err         error
	wroteHeader bool
	mu          sync.Mutex
}

func (w *timeoutWriter) Header() http.Header {
	// Header仅由Handler goroutine写入，并在done关闭后由外层读取。
	return w.header
}

func (w *timeoutWriter) WriteHeader(status int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.err != nil || w.wroteHeader {
		return
	}
	w.status = status
	w.wroteHeader = true
}

func (w *timeoutWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.err != nil {
		return 0, w.err
	}
	if !w.wroteHeader {
		w.status = http.StatusOK
		w.wroteHeader = true
	}
	return w.buffer.Write(data)
}

func (w *timeoutWriter) cancel(err error) {
	w.mu.Lock()
	w.err = err
	w.mu.Unlock()
}

func (w *timeoutWriter) writeResponse() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.err != nil {
		return
	}
	for key, values := range w.header {
		w.destination.Header()[key] = values
	}
	status := w.status
	if status == 0 {
		status = http.StatusOK
	}
	w.destination.WriteHeader(status)
	_, _ = w.destination.Write(w.buffer.Bytes())
}

func writeTimeoutBody(w http.ResponseWriter, err *errorx.Error) error {
	body := httpx.Body{
		Code: int(err.Code()),
		Msg:  err.Message(),
	}
	return json.NewEncoder(w).Encode(body)
}

func observeHandlerExit(
	log *zap.Logger,
	r *http.Request,
	done <-chan struct{},
	panicCh <-chan panicResult,
) {
	select {
	case <-done:
		return
	case result := <-panicCh:
		if result.value == http.ErrAbortHandler {
			log.Warn(
				"HTTP请求超时后被中止",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)
			return
		}
		log.Error(
			"HTTP请求超时后发生panic",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Any("panic", result.value),
			zap.ByteString("stack", result.stack),
		)
	}
}
