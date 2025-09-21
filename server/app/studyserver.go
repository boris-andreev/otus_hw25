package app

import (
	"context"
	"proto_api/pkg/grpc/v1/todo_api"
	"server/internal/model"
	"server/internal/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type studyServer struct {
	todoService *service.TodoService
	todo_api.UnimplementedStudyServiceServer
}

func (h *studyServer) CreateStudy(ctx context.Context, item *todo_api.CreateStudyRequest) (*emptypb.Empty, error) {
	h.todoService.CreateItem(&model.StudyItem{
		Topic: item.Topic,
	})

	return &emptypb.Empty{}, nil
}

func (h *studyServer) UpdateStudy(ctx context.Context, item *todo_api.UpdateStudyRequest) (*emptypb.Empty, error) {
	var id, err = primitive.ObjectIDFromHex(item.Id)
	if err != nil {
		return nil, nil
	}

	h.todoService.UpdateItem(&model.StudyItem{
		Id:    id,
		Topic: item.Topic,
	})

	return &emptypb.Empty{}, nil
}

func (h *studyServer) GetStudy(ctx context.Context, request *todo_api.GetStudyRequest) (*todo_api.GetStudyResponse, error) {
	item, err := h.todoService.GetStudyItem(request.Id)

	if err != nil {
		return nil, err
	}

	return &todo_api.GetStudyResponse{
		Id:    item.Id.Hex(),
		Topic: item.Topic,
	}, nil
}
func (h *studyServer) ListStudy(context.Context, *emptypb.Empty) (*todo_api.ListStudyResponse, error) {
	items, err := h.todoService.GetStudyItems()

	if err != nil {
		return nil, err
	}

	result := make([]*todo_api.GetStudyResponse, 0, len(items))

	for _, item := range items {
		result = append(result, &todo_api.GetStudyResponse{
			Id:    item.Id.Hex(),
			Topic: item.Topic,
		})
	}

	return &todo_api.ListStudyResponse{
		Result: result,
	}, nil
}

func (h *studyServer) DeleteStudy(ctx context.Context, request *todo_api.DeleteStudyRequest) (*emptypb.Empty, error) {
	h.todoService.DeleteStudyItem(request.Id)

	return &emptypb.Empty{}, nil
}

func newStudyServer(todoService *service.TodoService) *studyServer {
	return &studyServer{
		todoService: todoService,
	}
}
