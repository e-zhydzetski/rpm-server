package main

import (
	"context"
	"errors"
	"github.com/e-zhydzetski/rpm-server/internal/rpmserver"
	"github.com/go-chi/chi"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const addr = ":8080"

func main() {
	ctx := context.Background()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-c:
			return errors.New("graceful shutdown by " + sig.String())
		}
	})

	cfg := rpmserver.Config{
		AccessToken:    os.Getenv("ACCESS_TOKEN"),
		RepositoryPath: os.Getenv("REPO_PATH"),
	}

	r := chi.NewRouter()

	handler := rpmserver.NewHandler(cfg)
	r.Mount("/api", handler)
	r.Mount("/", http.FileServer(http.Dir(cfg.RepositoryPath)))

	server := &http.Server{
		Addr:    addr,
		Handler: chi.ServerBaseContext(ctx, r),
	}
	g.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(context.Background())
	})
	g.Go(func() error {
		log.Println("start listening at", addr, "...")
		return server.ListenAndServe()
	})

	err := g.Wait()
	log.Println("server stopped:", err)
}
