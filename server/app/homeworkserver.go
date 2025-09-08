package app

import (
	"context"
	"proto_api/pkg/grpc/v1/todo_api"
	"server/internal/model"
	"server/internal/service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type homeworkServer struct {
	todoService *service.TodoService
	todo_api.UnimplementedHomeworkServiceServer
}

func (h *homeworkServer) CreateHomework(ctx context.Context, item *todo_api.CreateHomeworkRequest) (*emptypb.Empty, error) {
	h.todoService.CreateItem(&model.HomeworkItem{
		Description: item.Description,
	})

	return &emptypb.Empty{}, nil
}

func (h *homeworkServer) UpdateHomework(ctx context.Context, item *todo_api.UpdateHomeworkRequest) (*emptypb.Empty, error) {
	h.todoService.UpdateItem(&model.HomeworkItem{
		Id:          int(item.Id),
		Description: item.Description,
	})

	return &emptypb.Empty{}, nil
}

func (h *homeworkServer) GetHomework(ctx context.Context, request *todo_api.GetHomeworkRequest) (*todo_api.GetHomeworkResponse, error) {
	item, err := h.todoService.GetHomeworkItem(int(request.Id))

	if err != nil {
		return nil, err
	}

	return &todo_api.GetHomeworkResponse{
		Id:          int64(item.Id),
		Description: item.Description,
	}, nil
}
func (h *homeworkServer) ListHomework(context.Context, *emptypb.Empty) (*todo_api.ListHomeworkResponse, error) {
	items, err := h.todoService.GetHomeworkItems()

	if err != nil {
		return nil, err
	}

	result := make([]*todo_api.GetHomeworkResponse, 0, len(items))

	for _, item := range items {
		result = append(result, &todo_api.GetHomeworkResponse{
			Id:          int64(item.Id),
			Description: item.Description,
		})
	}

	return &todo_api.ListHomeworkResponse{
		Result: result,
	}, nil
}

func (h *homeworkServer) DeleteHomework(ctx context.Context, request *todo_api.DeleteHomeworkRequest) (*emptypb.Empty, error) {
	h.todoService.DeleteHomeworkItem(int(request.Id))

	return &emptypb.Empty{}, nil
}

func newHomeworkServer(todoService *service.TodoService) *homeworkServer {
	return &homeworkServer{
		todoService: todoService,
	}
}
