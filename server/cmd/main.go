package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"server/app"
	"server/internal/postgresRepository"
	"server/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	todoService := service.NewTodoService(postgresrepository.NewTodoRepository(ctx, &wg), ctx, &wg)

	app := app.New(ctx, &wg, todoService)

	app.Start()

	wg.Wait()
	fmt.Println("\nGracefull shutdown is ok")
}
