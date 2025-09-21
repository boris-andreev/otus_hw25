package service

import (
	"context"
	"sync"
	"time"

	"server/internal/model"
)

type todoRepository interface {
	CreateItem(item model.Identifier)
	UpdateItem(item model.Identifier)
	DeleteHomeworkItem(id string) error
	DeleteStudyItem(id string) error
	DeleteWorkoutItem(id string) error
	GetHomeworkItem(id string) (*model.HomeworkItem, error)
	GetStudyItem(id string) (*model.StudyItem, error)
	GetWorkoutItem(id string) (*model.WorkoutItem, error)
	GetHomeworkItems() ([]*model.HomeworkItem, error)
	GetStudyItems() ([]*model.StudyItem, error)
	GetWorkoutItems() ([]*model.WorkoutItem, error)
	GetNewHomewors(timestamp time.Time) ([]*model.HomeworkItem, time.Time)
	GetNewStudies(timestamp time.Time) ([]*model.StudyItem, time.Time)
	GetNewWorkouts(timestamp time.Time) ([]*model.WorkoutItem, time.Time)
}

type TodoService struct {
	items      chan *operationItem
	repository todoRepository
	ctx        context.Context
	wg         *sync.WaitGroup
	logger     *logger
}

func (t *TodoService) CreateItem(item model.Identifier) {
	t.items <- &operationItem{item: item, operationType: add}
}

func (t *TodoService) UpdateItem(item model.Identifier) {
	t.items <- &operationItem{item: item, operationType: update}
}

func (t *TodoService) DeleteHomeworkItem(id string) error {
	return t.repository.DeleteHomeworkItem(id)
}

func (t *TodoService) DeleteStudyItem(id string) error {
	return t.repository.DeleteStudyItem(id)
}

func (t *TodoService) DeleteWorkoutItem(id string) error {
	return t.repository.DeleteWorkoutItem(id)
}

func (t *TodoService) GetHomeworkItem(id string) (*model.HomeworkItem, error) {
	return t.repository.GetHomeworkItem(id)
}

func (t *TodoService) GetStudyItem(id string) (*model.StudyItem, error) {
	return t.repository.GetStudyItem(id)
}

func (t *TodoService) GetWorkoutItem(id string) (*model.WorkoutItem, error) {
	return t.repository.GetWorkoutItem(id)
}

func (t *TodoService) GetHomeworkItems() ([]*model.HomeworkItem, error) {
	return t.repository.GetHomeworkItems()
}

func (t *TodoService) GetStudyItems() ([]*model.StudyItem, error) {
	return t.repository.GetStudyItems()
}

func (t *TodoService) GetWorkoutItems() ([]*model.WorkoutItem, error) {
	return t.repository.GetWorkoutItems()
}

func (t *TodoService) listenForItems() {
	var once sync.Once

	t.wg.Add(1)

	go func() {
		defer t.wg.Done()

		once.Do(func() {
			t.logger.Log()
		})

		for operationItem := range t.items {
			switch operationItem.operationType {
			case add:
				t.repository.CreateItem(operationItem.item)
			case update:
				t.repository.UpdateItem(operationItem.item)
			}
		}
	}()
}

func (t *TodoService) listenForFinish() {
	t.wg.Add(1)

	go func() {
		defer t.wg.Done()

		for {
			select {
			case <-t.ctx.Done():
				close(t.items)
				return
			}
		}
	}()
}

func NewTodoService(repo todoRepository, ctx context.Context, wg *sync.WaitGroup) *TodoService {
	res := &TodoService{
		items:      make(chan *operationItem),
		repository: repo,
		ctx:        ctx,
		wg:         wg,
		logger:     NewLogger(repo, ctx, wg),
	}

	res.listenForItems()
	res.listenForFinish()

	return res
}
