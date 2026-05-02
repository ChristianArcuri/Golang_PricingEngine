package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang_pricingengine/internal/config"
	"golang_pricingengine/internal/httpserver"
	"golang_pricingengine/internal/pricing"
	"golang_pricingengine/internal/repo"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.SQLServerConnectionString == "" {
		log.Fatal("SQLSERVER_CONNECTION_STRING is required")
	}
	rp, err := repo.Open(cfg.SQLServerConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer rp.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rp.Ping(ctx); err != nil {
		log.Printf("warning: database ping failed: %v", err)
	}

	svc := pricing.New(cfg, rp)
	srv := httpserver.New(svc)

	httpSrv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           srv.Router(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	_ = httpSrv.Shutdown(shutdownCtx)
}
