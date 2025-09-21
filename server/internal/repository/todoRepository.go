package repository

import (
	"context"
	"log"
	"server/internal/model"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	homeworks = "homeworks"
	workouts  = "workouts"
	studies   = "studies"
	db        = "todo"
)

type TodoRepositoryImpl struct {
	items  chan model.Identifier
	client *mongo.Client
	ctx    context.Context
	wg     *sync.WaitGroup
}

func (t *TodoRepositoryImpl) CreateItem(item model.Identifier) {
	switch item.(type) {
	case *model.HomeworkItem:
		appendItem(t.client, item, homeworks)
	case *model.StudyItem:
		appendItem(t.client, item, studies)
	case *model.WorkoutItem:
		appendItem(t.client, item, workouts)
	}
}

func appendItem(client *mongo.Client, item model.Identifier, entityName string) {
	collection := client.Database(db).Collection(entityName)
	item.SetId(primitive.NewObjectID())

	if _, err := collection.InsertOne(context.TODO(), item); err != nil {
		log.Fatal(err)
	}
}

func (t *TodoRepositoryImpl) UpdateItem(item model.Identifier) {
	switch item.(type) {
	case *model.HomeworkItem:
		saveItem(t.client, item, homeworks)
	case *model.StudyItem:
		saveItem(t.client, item, studies)
	case *model.WorkoutItem:
		saveItem(t.client, item, workouts)
	}
}

func saveItem(client *mongo.Client, item model.Identifier, entityName string) {
	collection := client.Database(db).Collection(entityName)

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": item.GetId()}, bson.M{"$set": item})
	if err != nil {
		log.Fatal(err)
	}
}

func (t *TodoRepositoryImpl) DeleteHomeworkItem(id string) error {
	return deleteItemById(t.client, id, homeworks)
}

func (t *TodoRepositoryImpl) DeleteStudyItem(id string) error {
	return deleteItemById(t.client, id, studies)
}

func (t *TodoRepositoryImpl) DeleteWorkoutItem(id string) error {
	return deleteItemById(t.client, id, workouts)
}

func deleteItemById(client *mongo.Client, entityName string, id string) error {
	collection := client.Database(db).Collection(entityName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoRepositoryImpl) GetHomeworkItem(id string) (*model.HomeworkItem, error) {
	return getItem[*model.HomeworkItem](t.client, homeworks, id)
}

func (t *TodoRepositoryImpl) GetStudyItem(id string) (*model.StudyItem, error) {
	return getItem[*model.StudyItem](t.client, studies, id)
}

func (t *TodoRepositoryImpl) GetWorkoutItem(id string) (*model.WorkoutItem, error) {
	return getItem[*model.WorkoutItem](t.client, workouts, id)
}

func getItem[T model.ItemWithId](client *mongo.Client, entityName string, id string) (T, error) {
	collection := client.Database(db).Collection(entityName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var t T

	err = collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TodoRepositoryImpl) GetHomeworkItems() ([]*model.HomeworkItem, error) {
	return getItems[*model.HomeworkItem](t.client, homeworks)
}

func (t *TodoRepositoryImpl) GetStudyItems() ([]*model.StudyItem, error) {
	return getItems[*model.StudyItem](t.client, studies)
}

func (t *TodoRepositoryImpl) GetWorkoutItems() ([]*model.WorkoutItem, error) {
	return getItems[*model.WorkoutItem](t.client, workouts)
}

func getItems[T model.ItemWithId](client *mongo.Client, entityName string) ([]T, error) {
	collection := client.Database(db).Collection(entityName)

	cursor, err := collection.Find(context.TODO(), primitive.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var res []T
	if err = cursor.All(context.TODO(), &res); err != nil {
		log.Fatal(err)
	}

	return res, nil
}

func (t *TodoRepositoryImpl) GetNewHomewors(timestamp time.Time) ([]*model.HomeworkItem, time.Time) {
	return getNewItems[*model.HomeworkItem](t.client, homeworks, timestamp)
}

func (t *TodoRepositoryImpl) GetNewStudies(timestamp time.Time) ([]*model.StudyItem, time.Time) {
	return getNewItems[*model.StudyItem](t.client, studies, timestamp)
}

func (t *TodoRepositoryImpl) GetNewWorkouts(timestamp time.Time) ([]*model.WorkoutItem, time.Time) {
	return getNewItems[*model.WorkoutItem](t.client, workouts, timestamp)
}

func getNewItems[T model.ItemWithId](client *mongo.Client, entityName string, timestamp time.Time) ([]T, time.Time) {
	collection := client.Database(db).Collection(entityName)

	filter := bson.M{
		"_id": bson.M{
			"$gte": primitive.NewObjectIDFromTimestamp(timestamp),
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Panic(err)
	}
	defer cursor.Close(context.TODO())

	var res []T
	if err = cursor.All(context.TODO(), &res); err != nil {
		log.Fatal(err)
	}

	return res, time.Now().UTC()
}

func NewTodoRepository(ctx context.Context, wg *sync.WaitGroup) *TodoRepositoryImpl {
	res := &TodoRepositoryImpl{
		ctx: ctx,
		wg:  wg,
	}

	res.initDb()
	res.listenForClosingDb()

	return res
}

func (t *TodoRepositoryImpl) initDb() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo-container:27017")
	ctx, cancel := context.WithTimeout(t.ctx, 10*time.Second)
	defer cancel()

	var err error
	t.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = t.client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
}

func (t *TodoRepositoryImpl) listenForClosingDb() {
	t.wg.Add(1)

	go func() {
		defer t.wg.Done()

		for {
			select {
			case <-t.ctx.Done():
				if t.client != nil {
					ctx, cancel := context.WithTimeout(t.ctx, 10*time.Second)
					defer cancel()

					if err := t.client.Disconnect(ctx); err != nil {
						log.Fatal(err)
					}
					log.Println("Disconnected from MongoDB.")
				}
				return
			}
		}
	}()
}
