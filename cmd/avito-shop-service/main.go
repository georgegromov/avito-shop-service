package main

import (
	"avito-shop-service/internal/app"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	a := app.New()

	go a.HttpServer.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.Close(ctx)
}
