package service

import (
	"context"
	"server/internal/model"
	"server/internal/service/mock"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var homeworkItem = &model.HomeworkItem{Description: "Description"}
var studyItem = &model.StudyItem{Topic: "Math"}
var workoutItem = &model.WorkoutItem{Target: "Biceps"}

func TestTodoService_CreateItem(t *testing.T) {
	type args struct {
		item model.Identifier
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Add homework",
			args: args{
				item: homeworkItem,
			},
			mock: func() {
				repositoryMock.EXPECT().CreateItem(homeworkItem)
			},
		},
		{
			name: "Add study",
			args: args{
				item: studyItem,
			},
			mock: func() {
				repositoryMock.EXPECT().CreateItem(studyItem)
			},
		},
		{
			name: "Add workout",
			args: args{
				item: workoutItem,
			},
			mock: func() {
				repositoryMock.EXPECT().CreateItem(workoutItem)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			tr.CreateItem(tt.args.item)

			cancel()

			wg.Wait()
		})
	}
}

func sharedMock(repositoryMock *mock.MocktodoRepository) {
	repositoryMock.EXPECT().GetNewHomewors(gomock.Any()).AnyTimes().Return(nil, time.Now().UTC())
	repositoryMock.EXPECT().GetNewStudies(gomock.Any()).AnyTimes().Return(nil, time.Now().UTC())
	repositoryMock.EXPECT().GetNewWorkouts(gomock.Any()).AnyTimes().Return(nil, time.Now().UTC())
}

func TestTodoService_UpdateItem(t *testing.T) {
	type args struct {
		item model.Identifier
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Update homework",
			args: args{
				item: homeworkItem,
			},
			mock: func() {
				repositoryMock.EXPECT().UpdateItem(homeworkItem)
			},
		},
		{
			name: "Update study",
			args: args{
				item: studyItem,
			},
			mock: func() {
				repositoryMock.EXPECT().UpdateItem(studyItem)
			},
		},
		{
			name: "Update workout",
			args: args{
				item: workoutItem,
			},
			mock: func() {
				repositoryMock.EXPECT().UpdateItem(workoutItem)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			tr.UpdateItem(tt.args.item)

			cancel()

			wg.Wait()
		})
	}
}

func TestTodoService_DeleteHomeworkItem(t *testing.T) {
	type args struct {
		id string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Delete homework",
			args: args{
				id: "10",
			},
			mock: func() {
				repositoryMock.EXPECT().DeleteHomeworkItem("10").Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			err := tr.DeleteHomeworkItem(tt.args.id)

			cancel()

			wg.Wait()

			assert.NoError(t, err)
		})
	}
}

func TestTodoService_DeleteStudyItem(t *testing.T) {
	type args struct {
		id string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Delete study",
			args: args{
				id: "10",
			},
			mock: func() {
				repositoryMock.EXPECT().DeleteStudyItem("10").Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			err := tr.DeleteStudyItem(tt.args.id)

			cancel()

			wg.Wait()

			assert.NoError(t, err)
		})
	}
}

func TestTodoService_DeleteWorkoutItem(t *testing.T) {
	type args struct {
		id string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Delete workout",
			args: args{
				id: "10",
			},
			mock: func() {
				repositoryMock.EXPECT().DeleteWorkoutItem("10").Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			err := tr.DeleteWorkoutItem(tt.args.id)

			cancel()

			wg.Wait()

			assert.NoError(t, err)
		})
	}
}

func TestTodoService_GetHomeworkItem(t *testing.T) {
	type args struct {
		id string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Get homework",
			args: args{
				id: "10",
			},
			mock: func() {
				repositoryMock.EXPECT().GetHomeworkItem("10").Return(homeworkItem, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			item, err := tr.GetHomeworkItem(tt.args.id)

			assert.Equal(t, homeworkItem, item)
			assert.NoError(t, err)

			cancel()

			wg.Wait()
		})
	}
}

func TestTodoService_GetStudyItem(t *testing.T) {
	type args struct {
		id string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Get study",
			args: args{
				id: "10",
			},
			mock: func() {
				repositoryMock.EXPECT().GetStudyItem("10").Return(studyItem, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			item, err := tr.GetStudyItem(tt.args.id)

			assert.Equal(t, studyItem, item)
			assert.NoError(t, err)

			cancel()

			wg.Wait()
		})
	}
}

func TestTodoService_GetWorkoutItem(t *testing.T) {
	type args struct {
		id string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Get workout",
			args: args{
				id: "10",
			},
			mock: func() {
				repositoryMock.EXPECT().GetWorkoutItem("10").Return(workoutItem, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			item, err := tr.GetWorkoutItem(tt.args.id)

			assert.Equal(t, workoutItem, item)
			assert.NoError(t, err)

			cancel()

			wg.Wait()
		})
	}
}

func TestTodoService_GetHomeworkItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		mock func()
	}{
		{
			name: "Get homeworks",
			mock: func() {
				repositoryMock.EXPECT().GetHomeworkItems().Return([]*model.HomeworkItem{homeworkItem}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			items, err := tr.GetHomeworkItems()

			assert.Equal(t, homeworkItem, items[0])
			assert.NoError(t, err)

			cancel()

			wg.Wait()
		})
	}
}

func TestTodoService_GetStudyItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		mock func()
	}{
		{
			name: "Get studies",
			mock: func() {
				repositoryMock.EXPECT().GetStudyItems().Return([]*model.StudyItem{studyItem}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			items, err := tr.GetStudyItems()

			assert.Equal(t, studyItem, items[0])
			assert.NoError(t, err)

			cancel()

			wg.Wait()
		})
	}
}

func TestTodoService_GetWorkoutItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock.NewMocktodoRepository(ctrl)

	sharedMock(repositoryMock)

	tests := []struct {
		name string
		mock func()
	}{
		{
			name: "Get workouts",
			mock: func() {
				repositoryMock.EXPECT().GetWorkoutItems().Return([]*model.WorkoutItem{workoutItem}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}

			tt.mock()

			tr := NewTodoService(repositoryMock, ctx, wg)
			items, err := tr.GetWorkoutItems()

			assert.Equal(t, workoutItem, items[0])
			assert.NoError(t, err)

			cancel()

			wg.Wait()
		})
	}
}
