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
		ListenAddr:    os.Getenv("LISTEN_ADDR"),
		AccessToken:   os.Getenv("ACCESS_TOKEN"),
		PushRepoPath:  os.Getenv("PUSH_PATH"),
		ReposRootPath: os.Getenv("REPOS_ROOT"),
	}
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":8080"
	}

	r := chi.NewRouter()

	handler := rpmserver.NewHandler(cfg)
	r.Mount("/api", handler)
	r.Mount("/repos",
		http.StripPrefix("/repos",
			rpmserver.NewFileServer(cfg.ReposRootPath),
		),
	)

	server := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: chi.ServerBaseContext(ctx, r),
	}
	g.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(context.Background())
	})
	g.Go(func() error {
		log.Println("start listening at", cfg.ListenAddr, "...")
		return server.ListenAndServe()
	})

	err := g.Wait()
	log.Println("server stopped:", err)
}
