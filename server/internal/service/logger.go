package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/internal/model"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type logger struct {
	repository todoRepository
	ctx        context.Context
	wg         *sync.WaitGroup
	rdb        *redis.Client
}

func (l *logger) Log() {
	l.wg.Add(1)

	go func() {
		defer l.wg.Done()

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		t := time.Now().UTC()
		homeworkLoggedTime := t
		studyLoggedTime := t
		workoutLoggedTime := t

		for {
			select {
			case <-ticker.C:
				func() {
					homeworkLoggedTime = logAddedItems(l.repository.GetNewHomewors, homeworkLoggedTime, l.rdb, l.ctx, "homework:%s")
					studyLoggedTime = logAddedItems(l.repository.GetNewStudies, studyLoggedTime, l.rdb, l.ctx, "study:%s")
					workoutLoggedTime = logAddedItems(l.repository.GetNewWorkouts, workoutLoggedTime, l.rdb, l.ctx, "workout:%s")
				}()
			case <-l.ctx.Done():
				return
			}
		}
	}()
}

func logAddedItems[T model.ItemWithId](
	getItems func(time.Time) ([]T, time.Time),
	loggedTime time.Time,
	rdb *redis.Client,
	ctx context.Context,
	keyFormat string) time.Time {

	items, timestamp := getItems(loggedTime)

	if len(items) > 0 {

		for _, item := range items {
			jsonData, err := json.Marshal(item)
			if err != nil {
				log.Panic(err)
			}

			err = rdb.Set(
				ctx, 
				fmt.Sprintf(keyFormat, item.GetId().Hex()), 
				jsonData, 
				time.Minute).Err()
			if err != nil {
				log.Panic(err)
			}
		}

		return timestamp
	}

	return loggedTime
}

func NewLogger(repo todoRepository, ctx context.Context, wg *sync.WaitGroup) *logger {

	return &logger{
		repository: repo,
		ctx:        ctx,
		wg:         wg,
		rdb: redis.NewClient(&redis.Options{
			Addr:     "redis-container:6379",
			Password: "",
			DB:       0,
		}),
	}
}
