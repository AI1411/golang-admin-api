package server

import (
	"context"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"

	"github.com/AI1411/golang-admin-api/logger"
)

const (
	port            = "8080"
	ShutdownTimeout = 30 * time.Second
)

func Run() {
	zapLogger, err := logger.NewLogger(true)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = zapLogger.Sync() }()

	zapLogger.Info("server start")

	router := Router()
	router.Use(gin.Middleware("golang-admin-api"))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Failure start server", zap.Error(err))
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Server shutdown:", zap.Error(err))
		panic(err)
	}
	zapLogger.Info("Success server shutdown")

	<-ctx.Done()

	zapLogger.Info("Server exiting")
}
