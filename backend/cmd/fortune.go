package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fortunefinger/backend/internal/startup"

	"go.uber.org/zap"
)

var configFile = flag.String("f", "../etc/fortune.yaml", "the config file")

func main() {
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := startup.NewApplication(ctx, *configFile)
	if err != nil {
		log.Fatalf("fortunefinger启动失败: %v", err)
	}

	runErr := application.Run(ctx)
	stop()
	closeErr := application.Close()
	exitErr := errors.Join(runErr, closeErr)
	if exitErr != nil {
		application.Logger().Error("fortunefinger service exited", zap.Error(exitErr))
	}
	_ = application.Logger().Sync()
	if exitErr != nil {
		os.Exit(1)
	}
}
