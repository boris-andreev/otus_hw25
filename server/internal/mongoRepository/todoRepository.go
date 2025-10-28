package mongorepository

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

type MongoTodoRepositoryImpl struct {
	items  chan model.Identifier
	client *mongo.Client
	ctx    context.Context
	wg     *sync.WaitGroup
}

func homeworkItemToBson(item *model.HomeworkItem) bson.M {
	doc := bson.M{
		"description": item.Description,
	}

	if item.Id != "" {
		objectID, err := primitive.ObjectIDFromHex(item.Id)
		if err == nil {
			doc["_id"] = objectID
		}
	} else {
		doc["_id"] = primitive.NewObjectID()
	}

	return doc
}

func homeworkItemFromBson(doc bson.M) *model.HomeworkItem {
	item := &model.HomeworkItem{}

	if id, ok := doc["_id"]; ok {
		if objectID, ok := id.(primitive.ObjectID); ok {
			item.Id = objectID.Hex()
		}
	}

	if description, ok := doc["description"]; ok {
		item.Description = description.(string)
	}

	return item
}

func studyItemToBson(item *model.StudyItem) bson.M {
	doc := bson.M{
		"topic": item.Topic,
	}

	if item.Id != "" {
		objectID, err := primitive.ObjectIDFromHex(item.Id)
		if err == nil {
			doc["_id"] = objectID
		}
	} else {
		doc["_id"] = primitive.NewObjectID()
	}

	return doc
}

func studyItemFromBson(doc bson.M) *model.StudyItem {
	item := &model.StudyItem{}

	if id, ok := doc["_id"]; ok {
		if objectID, ok := id.(primitive.ObjectID); ok {
			item.Id = objectID.Hex()
		}
	}

	if topic, ok := doc["topic"]; ok {
		item.Topic = topic.(string)
	}

	return item
}

func workoutItemToBson(item *model.WorkoutItem) bson.M {
	doc := bson.M{
		"target": item.Target,
	}

	if item.Id != "" {
		objectID, err := primitive.ObjectIDFromHex(item.Id)
		if err == nil {
			doc["_id"] = objectID
		}
	} else {
		doc["_id"] = primitive.NewObjectID()
	}

	return doc
}

func workoutItemFromBson(doc bson.M) *model.WorkoutItem {
	item := &model.WorkoutItem{}

	if id, ok := doc["_id"]; ok {
		if objectID, ok := id.(primitive.ObjectID); ok {
			item.Id = objectID.Hex()
		}
	}

	if target, ok := doc["target"]; ok {
		item.Target = target.(string)
	}

	return item
}

func (t *MongoTodoRepositoryImpl) CreateItem(item model.Identifier) {
	item.SetId(primitive.NewObjectID().Hex())

	switch item.(type) {
	case *model.HomeworkItem:
		appendItem(t.client, homeworkItemToBson(item.(*model.HomeworkItem)), homeworks)
	case *model.StudyItem:
		appendItem(t.client, studyItemToBson(item.(*model.StudyItem)), studies)
	case *model.WorkoutItem:
		appendItem(t.client, workoutItemToBson(item.(*model.WorkoutItem)), workouts)
	}
}

func appendItem(client *mongo.Client, item bson.M, entityName string) {
	collection := client.Database(db).Collection(entityName)

	if _, err := collection.InsertOne(context.TODO(), item); err != nil {
		log.Fatal(err)
	}
}

func (t *MongoTodoRepositoryImpl) UpdateItem(item model.Identifier) {
	switch item.(type) {
	case *model.HomeworkItem:
		saveItem(t.client, homeworkItemToBson(item.(*model.HomeworkItem)), homeworks)
	case *model.StudyItem:
		saveItem(t.client, studyItemToBson(item.(*model.StudyItem)), studies)
	case *model.WorkoutItem:
		saveItem(t.client, workoutItemToBson(item.(*model.WorkoutItem)), workouts)
	}
}

func saveItem(client *mongo.Client, item bson.M, entityName string) {
	collection := client.Database(db).Collection(entityName)

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": item["_id"]}, bson.M{"$set": item})
	if err != nil {
		log.Fatal(err)
	}
}

func (t *MongoTodoRepositoryImpl) DeleteHomeworkItem(id string) error {
	return deleteItemById(t.client, id, homeworks)
}

func (t *MongoTodoRepositoryImpl) DeleteStudyItem(id string) error {
	return deleteItemById(t.client, id, studies)
}

func (t *MongoTodoRepositoryImpl) DeleteWorkoutItem(id string) error {
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

func (t *MongoTodoRepositoryImpl) GetHomeworkItem(id string) (*model.HomeworkItem, error) {
	res, err := getItem[*model.HomeworkItem](t.client, homeworks, id)

	if err != nil {
		return nil, err
	}

	return homeworkItemFromBson(res), nil
}

func (t *MongoTodoRepositoryImpl) GetStudyItem(id string) (*model.StudyItem, error) {
	res, err := getItem[*model.StudyItem](t.client, studies, id)

	if err != nil {
		return nil, err
	}

	return studyItemFromBson(res), nil
}

func (t *MongoTodoRepositoryImpl) GetWorkoutItem(id string) (*model.WorkoutItem, error) {
	res, err := getItem[*model.WorkoutItem](t.client, workouts, id)

	if err != nil {
		return nil, err
	}

	return workoutItemFromBson(res), nil
}

func getItem[T model.ItemWithId](client *mongo.Client, entityName string, id string) (bson.M, error) {
	collection := client.Database(db).Collection(entityName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var t bson.M

	err = collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *MongoTodoRepositoryImpl) GetHomeworkItems() ([]*model.HomeworkItem, error) {
	var res []*model.HomeworkItem
	items, err := getItems[*model.HomeworkItem](t.client, homeworks)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		res = append(res, homeworkItemFromBson(item))
	}

	return res, nil
}

func (t *MongoTodoRepositoryImpl) GetStudyItems() ([]*model.StudyItem, error) {
	var res []*model.StudyItem
	items, err := getItems[*model.StudyItem](t.client, studies)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		res = append(res, studyItemFromBson(item))
	}

	return res, nil
}

func (t *MongoTodoRepositoryImpl) GetWorkoutItems() ([]*model.WorkoutItem, error) {
	var res []*model.WorkoutItem
	items, err := getItems[*model.WorkoutItem](t.client, workouts)

	if err != nil {
		return nil, err
	}

	for _, item := range items {
		res = append(res, workoutItemFromBson(item))
	}

	return res, nil
}

func getItems[T model.ItemWithId](client *mongo.Client, entityName string) ([]bson.M, error) {
	collection := client.Database(db).Collection(entityName)

	cursor, err := collection.Find(context.TODO(), primitive.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var res []bson.M
	if err = cursor.All(context.TODO(), &res); err != nil {
		log.Fatal(err)
	}

	return res, nil
}

func (t *MongoTodoRepositoryImpl) GetNewHomewors(timestamp time.Time) ([]*model.HomeworkItem, time.Time) {
	var res []*model.HomeworkItem
	items, ts := getNewItems[*model.HomeworkItem](t.client, homeworks, timestamp)

	for _, item := range items {
		res = append(res, homeworkItemFromBson(item))
	}

	return res, ts
}

func (t *MongoTodoRepositoryImpl) GetNewStudies(timestamp time.Time) ([]*model.StudyItem, time.Time) {
	var res []*model.StudyItem
	items, ts := getNewItems[*model.StudyItem](t.client, studies, timestamp)

	for _, item := range items {
		res = append(res, studyItemFromBson(item))
	}

	return res, ts
}

func (t *MongoTodoRepositoryImpl) GetNewWorkouts(timestamp time.Time) ([]*model.WorkoutItem, time.Time) {
	var res []*model.WorkoutItem
	items, ts := getNewItems[*model.WorkoutItem](t.client, workouts, timestamp)

	for _, item := range items {
		res = append(res, workoutItemFromBson(item))
	}

	return res, ts
}

func getNewItems[T model.ItemWithId](client *mongo.Client, entityName string, timestamp time.Time) ([]bson.M, time.Time) {
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

	var res []bson.M
	if err = cursor.All(context.TODO(), &res); err != nil {
		log.Fatal(err)
	}

	return res, time.Now().UTC()
}

func NewTodoRepository(ctx context.Context, wg *sync.WaitGroup) *MongoTodoRepositoryImpl {
	res := &MongoTodoRepositoryImpl{
		ctx: ctx,
		wg:  wg,
	}

	res.initDb()
	res.listenForClosingDb()

	return res
}

func (t *MongoTodoRepositoryImpl) initDb() {
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

func (t *MongoTodoRepositoryImpl) listenForClosingDb() {
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
