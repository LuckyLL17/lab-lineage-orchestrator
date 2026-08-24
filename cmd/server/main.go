package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/local/lab-lineage-orchestrator/internal/app"
	"github.com/local/lab-lineage-orchestrator/internal/httpapi"
	"github.com/local/lab-lineage-orchestrator/internal/jobs"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func main() {
	config := platform.LoadConfig()
	logger := platform.NewLogger()
	service := app.NewService(platform.RealClock{}, logger)
	_ = service.Bootstrap("system")
	router := httpapi.NewRouter(service, logger)
	runner := jobs.NewRunner(service, logger, config.Tick)

	server := &http.Server{Addr: config.Addr, Handler: router.Handler()}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	runner.Start(ctx)
	logger.Info("service starting", "addr", config.Addr, "domain", service.Describe()["domain"])

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), config.Shutdown)
		defer cancel()
		runner.Stop()
		_ = server.Shutdown(shutdownCtx)
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
