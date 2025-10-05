package app

import (
	"context"
	"proto_api/pkg/grpc/v1/todo_api"
	"server/internal/model"
	"server/internal/service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type workoutServer struct {
	todoService *service.TodoService
	todo_api.UnimplementedWorkoutServiceServer
}

func (h *workoutServer) CreateWorkout(ctx context.Context, item *todo_api.CreateWorkoutRequest) (*emptypb.Empty, error) {
	h.todoService.CreateItem(&model.WorkoutItem{
		Target: item.Target,
	})

	return &emptypb.Empty{}, nil
}

func (h *workoutServer) UpdateWorkout(ctx context.Context, item *todo_api.UpdateWorkoutRequest) (*emptypb.Empty, error) {
	h.todoService.UpdateItem(&model.WorkoutItem{
		Id:     item.Id,
		Target: item.Target,
	})

	return &emptypb.Empty{}, nil
}

func (h *workoutServer) GetWorkout(ctx context.Context, request *todo_api.GetWorkoutRequest) (*todo_api.GetWorkoutResponse, error) {
	item, err := h.todoService.GetWorkoutItem(request.Id)

	if err != nil {
		return nil, err
	}

	return &todo_api.GetWorkoutResponse{
		Id:     item.Id,
		Target: item.Target,
	}, nil
}
func (h *workoutServer) ListWorkout(context.Context, *emptypb.Empty) (*todo_api.ListWorkoutResponse, error) {
	items, err := h.todoService.GetWorkoutItems()

	if err != nil {
		return nil, err
	}

	result := make([]*todo_api.GetWorkoutResponse, 0, len(items))

	for _, item := range items {
		result = append(result, &todo_api.GetWorkoutResponse{
			Id:     item.Id,
			Target: item.Target,
		})
	}

	return &todo_api.ListWorkoutResponse{
		Result: result,
	}, nil
}

func (h *workoutServer) DeleteWorkout(ctx context.Context, request *todo_api.DeleteWorkoutRequest) (*emptypb.Empty, error) {
	h.todoService.DeleteWorkoutItem(request.Id)

	return &emptypb.Empty{}, nil
}

func newWorkoutServer(todoService *service.TodoService) *workoutServer {
	return &workoutServer{
		todoService: todoService,
	}
}
