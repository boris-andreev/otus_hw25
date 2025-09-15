package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type logger struct {
	repository todoRepository
	ctx        context.Context
	wg         *sync.WaitGroup
}

func (l *logger) Log() {
	l.wg.Add(1)

	go func() {
		defer l.wg.Done()

		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		lastHomeworkItemIdLogged := l.repository.GetLasttHomeworkItemId()
		lastStudyItemIdLogged := l.repository.GetLastStudyItemId()
		lastWorkoutItemIdLogged := l.repository.GetLastWorkoutItemId()

		for {
			select {
			case <-ticker.C:
				func() {
					lastHomeworkItemIdLogged = logAddedItems(lastHomeworkItemIdLogged, "Homeworks were added:", l.repository.GetNewHomewors)
					lastStudyItemIdLogged = logAddedItems(lastStudyItemIdLogged, "Studies were added:", l.repository.GetNewStudies)
					lastWorkoutItemIdLogged = logAddedItems(lastWorkoutItemIdLogged, "Workouts were added:", l.repository.GetNewWorkouts)
				}()
			case <-l.ctx.Done():
				return
			}
		}
	}()
}

func logAddedItems[T any](lastItemIdLogged int, message string, getItems func(int) (int, []*T)) int {
	lastItemId, items := getItems(lastItemIdLogged)

	if lastItemId > lastItemIdLogged {
		sb := strings.Builder{}
		sb.Write([]byte(message))
		sb.WriteRune('[')

		for i, item := range items {
			sb.WriteString(fmt.Sprintf("%v", *item))
			if i < len(items)-1 {
				sb.WriteString(", ")
			}
		}

		sb.WriteRune(']')
		log.Print(sb.String())
	}

	return lastItemId
}

func NewLogger(repo todoRepository, ctx context.Context, wg *sync.WaitGroup) *logger {
	return &logger{
		repository: repo,
		ctx:        ctx,
		wg:         wg,
	}
}
