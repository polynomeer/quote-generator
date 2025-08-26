package main

import (
	"context"
	"github.com/polynomeer/quote-generator/internal/config"
	"github.com/polynomeer/quote-generator/internal/generator"
	"github.com/polynomeer/quote-generator/internal/ingress"
	"github.com/polynomeer/quote-generator/internal/metrics"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfgPath := os.Getenv("QG_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/config.yaml"
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 생성기 초기화 (스텁: 1Hz 기본)
	gen := generator.New(generator.Options{
		Seed:    cfg.Seed,
		Symbols: cfg.Symbols,
		Hz:      cfg.Hz,
	})
	gen.Start()
	defer gen.Stop()

	mux := http.NewServeMux()
	ingress.RegisterHTTP(mux, gen)
	metrics.Register(mux)

	srv := &http.Server{
		Addr:           cfg.HTTP.Addr,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("quote-generator listening on %s", cfg.HTTP.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
