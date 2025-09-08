package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"hw25/internal/repository"
	"hw25/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	service.NewTodoServise(repository.NewTodoRepository(), ctx, &wg)

	wg.Wait()
	fmt.Println("\nGracefull shutdown is ok")
}
