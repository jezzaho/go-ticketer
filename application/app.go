package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
	srvcfg struct {
		Addr string
	}
}

func New() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{}),
		srvcfg: struct {
			Addr string
		}{Addr: ":3030"},
	}
	app.loadRoutes()
	return app
}

func (app *App) Start(ctx context.Context) error {
	server := http.Server{
		Addr:    app.srvcfg.Addr,
		Handler: app.router,
	}

	err := app.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis; %w", err)
	}

	defer func() {
		if err := app.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting server...")

	// channels for graceful shutdown
	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %s", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		fmt.Println("\nShutting down the server...")
		return server.Shutdown(timeout)
	}
}
