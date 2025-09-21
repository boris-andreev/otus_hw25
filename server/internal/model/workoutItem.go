package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutItem struct {
	Id     primitive.ObjectID `json:"id"`
	Target string             `json:"target" binding:"required"`
}

func (w *WorkoutItem) GetId() primitive.ObjectID {
	return w.Id
}

func (w *WorkoutItem) SetId(id primitive.ObjectID) {
	w.Id = id
}
