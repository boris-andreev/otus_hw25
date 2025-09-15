package service

import (
	"context"
	"sync"

	"server/internal/model"
)

type todoRepository interface {
	CreateItem(item model.Identifier)
	UpdateItem(item model.Identifier)
	GetLasttHomeworkItemId() int
	GetLastStudyItemId() int
	GetLastWorkoutItemId() int
	DeleteHomeworkItem(id int) error
	DeleteStudyItem(id int) error
	DeleteWorkoutItem(id int) error
	GetHomeworkItem(id int) (*model.HomeworkItem, error)
	GetStudyItem(id int) (*model.StudyItem, error)
	GetWorkoutItem(id int) (*model.WorkoutItem, error)
	GetHomeworkItems() ([]*model.HomeworkItem, error)
	GetStudyItems() ([]*model.StudyItem, error)
	GetWorkoutItems() ([]*model.WorkoutItem, error)
	GetNewHomewors(lastHomeworkItemId int) (int, []*model.HomeworkItem)
	GetNewStudies(lastStudyItemId int) (int, []*model.StudyItem)
	GetNewWorkouts(lastWorkoutItemId int) (int, []*model.WorkoutItem)
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

func (t *TodoService) DeleteHomeworkItem(id int) error {
	return t.repository.DeleteHomeworkItem(id)
}

func (t *TodoService) DeleteStudyItem(id int) error {
	return t.repository.DeleteStudyItem(id)
}

func (t *TodoService) DeleteWorkoutItem(id int) error {
	return t.repository.DeleteWorkoutItem(id)
}

func (t *TodoService) GetHomeworkItem(id int) (*model.HomeworkItem, error) {
	return t.repository.GetHomeworkItem(id)
}

func (t *TodoService) GetStudyItem(id int) (*model.StudyItem, error) {
	return t.repository.GetStudyItem(id)
}

func (t *TodoService) GetWorkoutItem(id int) (*model.WorkoutItem, error) {
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
