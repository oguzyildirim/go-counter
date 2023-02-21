package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Oguzyildirim/go-counter/internal/counter/rest"
	"github.com/Oguzyildirim/go-counter/internal/counter/service"
	cdb "github.com/Oguzyildirim/go-counter/internal/counter/storage"
	"github.com/Oguzyildirim/go-counter/internal/limiter/middleware"
	ldb "github.com/Oguzyildirim/go-counter/internal/limiter/storage"
)

func main() {
	var counterdb, limiterdb, address string
	flag.StringVar(&counterdb, "counterdb", "counterdb", "path to counterdb")
	flag.StringVar(&limiterdb, "limiterdb", "limiterdb", "path to limiterdb")
	flag.StringVar(&address, "address", ":9234", "HTTP Server Address")
	flag.Parse()

	if _, err := os.Stat(counterdb); os.IsNotExist(err) {
		_, err := os.Create(counterdb)
		if err != nil {
			log.Fatalf("Couldn't run: %s", err)
		}
	}

	if _, err := os.Stat(limiterdb); os.IsNotExist(err) {
		_, err := os.Create(limiterdb)
		if err != nil {
			log.Fatalf("Couldn't run: %s", err)
		}
	}

	errC, err := run(address, counterdb, limiterdb)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(address, counterdb string, limiterdb string) (<-chan error, error) {

	errC := make(chan error, 1)

	srv := newServer(address, counterdb, limiterdb)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		log.Printf("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		log.Printf("Shutdown completed")
	}()

	go func() {
		log.Printf("Listening and serving %s", address)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	return errC, nil

}

func newServer(address string, counterdb string, limiterdb string) *http.Server {
	r := http.NewServeMux()

	// limiter
	limiterRepo := ldb.NewLimiter(limiterdb)
	rateLimiter := middleware.NewRateLimitMiddleware(limiterRepo)

	// counter
	counterRepo := cdb.NewCounter(counterdb)
	counterSvc := service.NewCounter(counterRepo)

	rest.NewCounterHandler(counterSvc, rateLimiter).Register(r)

	return &http.Server{
		Handler:           r,
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
}
