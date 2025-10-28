package postgresrepository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"server/internal/model"
	"sync"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

type PostgresTodoRepositoryImpl struct {
	items chan model.Identifier
	db    *sql.DB
	ctx   context.Context
	wg    *sync.WaitGroup
}

func (t *PostgresTodoRepositoryImpl) CreateItem(item model.Identifier) {
	switch item.(type) {
	case *model.HomeworkItem:
		query := `
		INSERT INTO homeworks (description) 
		VALUES ($1)`
		t.execQuery(query, item.(*model.HomeworkItem).Description)
	case *model.StudyItem:
		query := `
		INSERT INTO studies (topic) 
		VALUES ($1)`
		t.execQuery(query, item.(*model.StudyItem).Topic)
	case *model.WorkoutItem:
		query := `
		INSERT INTO workouts (target) 
		VALUES ($1)`
		t.execQuery(query, item.(*model.WorkoutItem).Target)
	}
}

func (t *PostgresTodoRepositoryImpl) UpdateItem(item model.Identifier) {
	switch item.(type) {
	case *model.HomeworkItem:
		query := `
		UPDATE homeworks 
		SET description = $1
		WHERE id = $2`
		t.execQuery(query, item.(*model.HomeworkItem).Description, item.GetId())
	case *model.StudyItem:
		query := `
		UPDATE studies
		SET topic = $1
		WHERE id = $2`
		t.execQuery(query, item.(*model.StudyItem).Topic, item.GetId())
	case *model.WorkoutItem:
		query := `
		UPDATE workouts
		SET target = $1
		WHERE id = $2`
		t.execQuery(query, item.(*model.WorkoutItem).Target, item.GetId())
	}
}

func (t *PostgresTodoRepositoryImpl) execQuery(query string, args ...any) error {
	_, err := t.db.ExecContext(t.ctx, query, args...)

	if err != nil {
		log.Println(err)
	}

	return err
}

func (t *PostgresTodoRepositoryImpl) DeleteHomeworkItem(id string) error {
	query := `
	DELETE FROM homeworks 
	WHERE id = $1`

	return t.execQuery(query, id)
}

func (t *PostgresTodoRepositoryImpl) DeleteStudyItem(id string) error {
	query := `
	DELETE from studies 
	WHERE id = $1`

	return t.execQuery(query, id)
}

func (t *PostgresTodoRepositoryImpl) DeleteWorkoutItem(id string) error {
	query := `
	DELETE from workouts 
	WHERE id = $1`

	return t.execQuery(query, id)
}

func (t *PostgresTodoRepositoryImpl) GetHomeworkItem(id string) (*model.HomeworkItem, error) {
	res := &model.HomeworkItem{}

	query := "SELECT id, description FROM studies WHERE id = $1"

	err := t.db.QueryRowContext(t.ctx, query, id).Scan(&res.Id, &res.Description)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *PostgresTodoRepositoryImpl) GetStudyItem(id string) (*model.StudyItem, error) {
	res := &model.StudyItem{}

	query := "SELECT id, topic FROM studies WHERE id = $1"

	err := t.db.QueryRowContext(t.ctx, query, id).Scan(&res.Id, &res.Topic)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *PostgresTodoRepositoryImpl) GetWorkoutItem(id string) (*model.WorkoutItem, error) {
	res := &model.WorkoutItem{}

	query := "SELECT id, target FROM workouts WHERE id = $1"

	err := t.db.QueryRowContext(t.ctx, query, id).Scan(&res.Id, &res.Target)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *PostgresTodoRepositoryImpl) GetHomeworkItems() ([]*model.HomeworkItem, error) {
	query := `
    SELECT id, description 
    FROM homeworks
    ORDER BY id`

	return t.getHomeworkItems(query)
}

func (t *PostgresTodoRepositoryImpl) getHomeworkItems(query string, args ...any) ([]*model.HomeworkItem, error) {
	var res []*model.HomeworkItem

	rows, err := t.db.QueryContext(t.ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &model.HomeworkItem{}
		err := rows.Scan(&item.Id, &item.Description)

		if err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	return res, nil
}

func (t *PostgresTodoRepositoryImpl) GetStudyItems() ([]*model.StudyItem, error) {
	query := `
    SELECT id, topic 
    FROM studies
    ORDER BY id`

	return t.getStudyItems(query)
}

func (t *PostgresTodoRepositoryImpl) getStudyItems(query string, args ...any) ([]*model.StudyItem, error) {
	var res []*model.StudyItem

	rows, err := t.db.QueryContext(t.ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &model.StudyItem{}
		err := rows.Scan(&item.Id, &item.Topic)

		if err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	return res, nil
}

func (t *PostgresTodoRepositoryImpl) GetWorkoutItems() ([]*model.WorkoutItem, error) {
	query := `
    SELECT id, target 
    FROM workouts
    ORDER BY id`

	return t.getWorkoutItems(query)
}

func (t *PostgresTodoRepositoryImpl) getWorkoutItems(query string, args ...any) ([]*model.WorkoutItem, error) {
	var res []*model.WorkoutItem

	rows, err := t.db.QueryContext(t.ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &model.WorkoutItem{}
		err := rows.Scan(&item.Id, &item.Target)

		if err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	return res, nil
}

func (t *PostgresTodoRepositoryImpl) GetNewHomewors(timestamp time.Time) ([]*model.HomeworkItem, time.Time) {
	query := `
    SELECT id, description 
    FROM homeworks
	WHERE uuid_extract_timestamp(id) >= $1
    ORDER BY id`

	res, err := t.getHomeworkItems(query, timestamp)

	if err != nil {
		log.Fatal(err)
	}

	return res, time.Now().UTC()
}

func (t *PostgresTodoRepositoryImpl) GetNewStudies(timestamp time.Time) ([]*model.StudyItem, time.Time) {
	query := `
    SELECT id, topic 
    FROM studies
	WHERE uuid_extract_timestamp(id) >= $1
    ORDER BY id`

	res, err := t.getStudyItems(query, timestamp)

	if err != nil {
		log.Fatal(err)
	}

	return res, time.Now().UTC()
}

func (t *PostgresTodoRepositoryImpl) GetNewWorkouts(timestamp time.Time) ([]*model.WorkoutItem, time.Time) {
	query := `
    SELECT id, target
    FROM workouts
	WHERE uuid_extract_timestamp(id) >= $1
    ORDER BY id`

	res, err := t.getWorkoutItems(query, timestamp)

	if err != nil {
		log.Fatal(err)
	}

	return res, time.Now().UTC()
}

func NewTodoRepository(ctx context.Context, wg *sync.WaitGroup) *PostgresTodoRepositoryImpl {
	res := &PostgresTodoRepositoryImpl{
		ctx: ctx,
		wg:  wg,
	}

	res.initDb()
	res.listenForClosingDb()

	return res
}

func (t *PostgresTodoRepositoryImpl) initDb() {
	var err error

	t.db, err = sql.Open("pgx", os.Getenv("POSTGRES_CONN"))

	if err != nil {
		log.Fatalf("failed to load driver: %v", err)
	}

	err = t.db.PingContext(t.ctx)
	if err != nil {
		log.Fatalf("failed to connect to db: %v, %v", err, os.Getenv("POSTGRES_CONN"))
	}

	log.Println("Connected to Postgres!")
}

func (t *PostgresTodoRepositoryImpl) listenForClosingDb() {
	t.wg.Add(1)

	go func() {
		defer t.wg.Done()

		for {
			select {
			case <-t.ctx.Done():
				if t.db != nil {
					if err := t.db.Close(); err != nil {
						log.Fatal(err)
					}

					log.Println("Disconnected from MongoDB.")
				}
				return
			}
		}
	}()
}
