package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"proto_api/pkg/grpc/v1/todo_api"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"server/internal/service"
)

type App struct {
	todoService *service.TodoService
	server      *grpc.Server
	ctx         context.Context
	wg          *sync.WaitGroup
}

func New(ctx context.Context, wg *sync.WaitGroup, todoService *service.TodoService) *App {

	return &App{
		server:      grpc.NewServer(),
		ctx:         ctx,
		wg:          wg,
		todoService: todoService,
	}
}

func (a *App) Start() {
	go func() {
		l, err := net.Listen("tcp", "localhost:5001")

		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		todo_api.RegisterHomeworkServiceServer(a.server, newHomeworkServer(a.todoService))
		todo_api.RegisterStudyServiceServer(a.server, newStudyServer(a.todoService))
		todo_api.RegisterWorkoutServiceServer(a.server, newWorkoutServer(a.todoService))

		reflection.Register(a.server)

		log.Println("Listen at localhost:5001")

		if err := a.server.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("Server failed to start: %v", err))
		}
	}()

	a.listenForFinish()
}

func (a *App) listenForFinish() {
	a.wg.Add(1)

	go func() {
		defer a.wg.Done()

		for {
			select {
			case <-a.ctx.Done():
				fmt.Println("Shutting down server...")

				a.server.GracefulStop()
				log.Println("Server exited")
				return
			}
		}
	}()
}
