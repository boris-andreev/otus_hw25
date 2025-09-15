package repository

import (
	"log"
	"sort"
	"sync"

	"server/internal/model"
)

const (
	homeworksJson = "./homeworks.json"
	studiesJson   = "./studies.json"
	workoutsJson  = "./workouts.json"
)

type TodoRepository struct {
	homeworks []*model.HomeworkItem
	studies   []*model.StudyItem
	workouts  []*model.WorkoutItem

	homeworksMutex *sync.RWMutex
	studiesMutex   *sync.RWMutex
	workoutsMutex  *sync.RWMutex

	items chan model.Identifier
}

func (t *TodoRepository) CreateItem(item model.Identifier) {
	switch item.(type) {
	case *model.HomeworkItem:
		appendItem[*model.HomeworkItem](&t.homeworks, item.(*model.HomeworkItem), t.homeworksMutex, homeworksJson)
	case *model.StudyItem:
		appendItem[*model.StudyItem](&t.studies, item.(*model.StudyItem), t.studiesMutex, studiesJson)
	case *model.WorkoutItem:
		appendItem[*model.WorkoutItem](&t.workouts, item.(*model.WorkoutItem), t.workoutsMutex, workoutsJson)
	}
}

func appendItem[T model.ItemWithId](
	slice *[]T,
	item T,
	mutex *sync.RWMutex,
	fileName string) {

	mutex.Lock()
	defer mutex.Unlock()

	item.SetId(len(*slice) + 1)

	*slice = append(*slice, item)
	err := appendToFile(fileName, item)
	if err != nil {
		log.Panic(err)
	}
}

func (t *TodoRepository) UpdateItem(item model.Identifier) {
	switch item.(type) {
	case *model.HomeworkItem:
		saveItem[*model.HomeworkItem](t.homeworks, item.(*model.HomeworkItem), t.homeworksMutex, homeworksJson)
	case *model.StudyItem:
		saveItem[*model.StudyItem](t.studies, item.(*model.StudyItem), t.studiesMutex, studiesJson)
	case *model.WorkoutItem:
		saveItem[*model.WorkoutItem](t.workouts, item.(*model.WorkoutItem), t.workoutsMutex, workoutsJson)
	}
}

func saveItem[T model.ItemWithId](
	slice []T,
	item T,
	mutex *sync.RWMutex,
	fileName string) {

	mutex.Lock()
	defer mutex.Unlock()

	found, id := searchItemIndexById(slice, item.GetId())
	if !found {
		log.Printf("Could not find item with id: %d", item.GetId())
	}

	slice[id] = item

	if err := overwriteFile(fileName, slice); err != nil {
		log.Panic(err)
	}
}

func (t *TodoRepository) GetLasttHomeworkItemId() int {
	t.homeworksMutex.RLock()
	defer t.homeworksMutex.RUnlock()

	length := len(t.homeworks)
	if length == 0 {
		return 0
	}

	return t.homeworks[length-1].Id
}

func (t *TodoRepository) GetLastStudyItemId() int {
	t.studiesMutex.RLock()
	defer t.studiesMutex.RUnlock()

	length := len(t.studies)
	if length == 0 {
		return 0
	}

	return t.studies[length-1].Id
}

func (t *TodoRepository) GetLastWorkoutItemId() int {
	t.workoutsMutex.RLock()
	defer t.workoutsMutex.RUnlock()

	length := len(t.workouts)
	if length == 0 {
		return 0
	}

	return t.workouts[length-1].Id
}

func (t *TodoRepository) DeleteHomeworkItem(id int) error {
	t.homeworksMutex.Lock()
	defer t.homeworksMutex.Unlock()

	if deleteItemById(&t.homeworks, id) {
		return overwriteFile(homeworksJson, t.homeworks)
	}

	return nil
}

func (t *TodoRepository) DeleteStudyItem(id int) error {
	t.studiesMutex.Lock()
	defer t.studiesMutex.Unlock()

	if deleteItemById(&t.studies, id) {
		return overwriteFile(studiesJson, t.studies)
	}

	return nil
}

func (t *TodoRepository) DeleteWorkoutItem(id int) error {
	t.workoutsMutex.Lock()
	defer t.workoutsMutex.Unlock()

	if deleteItemById(&t.workouts, id) {
		return overwriteFile(workoutsJson, t.workouts)
	}

	return nil
}

func (t *TodoRepository) GetHomeworkItem(id int) (*model.HomeworkItem, error) {
	t.homeworksMutex.RLock()
	defer t.homeworksMutex.RUnlock()

	found, id := searchItemIndexById(t.homeworks, id)

	if !found {
		return nil, nil
	}

	return t.homeworks[id], nil
}

func (t *TodoRepository) GetStudyItem(id int) (*model.StudyItem, error) {
	t.studiesMutex.Lock()
	defer t.studiesMutex.Unlock()

	found, id := searchItemIndexById(t.studies, id)

	if !found {
		return nil, nil
	}

	return t.studies[id], nil
}

func (t *TodoRepository) GetWorkoutItem(id int) (*model.WorkoutItem, error) {
	t.workoutsMutex.Lock()
	defer t.workoutsMutex.Unlock()

	found, id := searchItemIndexById(t.workouts, id)

	if !found {
		return nil, nil
	}

	return t.workouts[id], nil
}

func (t *TodoRepository) GetHomeworkItems() ([]*model.HomeworkItem, error) {
	t.homeworksMutex.RLock()
	defer t.homeworksMutex.RUnlock()

	return t.homeworks, nil
}

func (t *TodoRepository) GetStudyItems() ([]*model.StudyItem, error) {
	t.studiesMutex.Lock()
	defer t.studiesMutex.Unlock()

	return t.studies, nil
}

func (t *TodoRepository) GetWorkoutItems() ([]*model.WorkoutItem, error) {
	t.workoutsMutex.Lock()
	defer t.workoutsMutex.Unlock()

	return t.workouts, nil
}

func deleteItemById[T model.ItemWithId](slice *[]T, id int) bool {
	found, i := searchItemIndexById(*slice, id)
	if found {
		*slice = append((*slice)[:i], (*slice)[i+1:]...)
	}

	return found
}

func searchItemIndexById[T model.ItemWithId](slice []T, id int) (bool, int) {
	i := sort.Search(len(slice), func(i int) bool {
		return slice[i].GetId() >= id
	})

	res := i < len(slice) && slice[i].GetId() == id

	return res, i
}

func (t *TodoRepository) GetNewHomewors(lastHomeworkItemId int) (int, []*model.HomeworkItem) {
	return getNewItems(lastHomeworkItemId, t.homeworksMutex, t.homeworks)
}

func (t *TodoRepository) GetNewStudies(lastStudyItemId int) (int, []*model.StudyItem) {
	return getNewItems(lastStudyItemId, t.studiesMutex, t.studies)
}

func (t *TodoRepository) GetNewWorkouts(lastWorkoutItemId int) (int, []*model.WorkoutItem) {
	return getNewItems(lastWorkoutItemId, t.workoutsMutex, t.workouts)
}

func getNewItems[T model.ItemWithId](
	lastItemId int,
	mu *sync.RWMutex,
	slice []T) (int, []T) {

	mu.RLock()
	defer mu.RUnlock()

	doesNextItemExist, startIndex, lastId, length := getNextItemIndexById(lastItemId, slice)

	if !doesNextItemExist {
		return lastItemId, []T{}
	}

	sliceCopy := make([]T, length-startIndex)
	copy(sliceCopy, slice[startIndex:length])

	return lastId, sliceCopy
}

func getNextItemIndexById[T model.ItemWithId](itemId int, slice []T) (doesNextItemExist bool, itemIndex int, lastId int, length int) {

	length = len(slice)

	if length == 0 {
		return false, 0, 0, length
	}

	startIndex := length - 1

	lastId = slice[startIndex].GetId()

	if lastId == itemId {
		return false, 0, 0, length
	}

	for startIndex >= 0 {
		id := slice[startIndex].GetId()
		if id <= itemId {
			break
		}
		itemIndex = startIndex

		startIndex--
	}

	return true, itemIndex, lastId, length
}

func NewTodoRepository() *TodoRepository {
	result := &TodoRepository{
		homeworksMutex: &sync.RWMutex{},
		studiesMutex:   &sync.RWMutex{},
		workoutsMutex:  &sync.RWMutex{},
	}
	var err error

	result.homeworksMutex.Lock()
	result.homeworks, err = readFromFile[model.HomeworkItem](homeworksJson)
	if err != nil {
		panic(err)
	}
	result.homeworksMutex.Unlock()

	result.studiesMutex.Lock()
	result.studies, err = readFromFile[model.StudyItem](studiesJson)
	if err != nil {
		panic(err)
	}
	result.studiesMutex.Unlock()

	result.workoutsMutex.Lock()
	result.workouts, err = readFromFile[model.WorkoutItem](workoutsJson)
	if err != nil {
		panic(err)
	}
	result.workoutsMutex.Unlock()

	return result
}
