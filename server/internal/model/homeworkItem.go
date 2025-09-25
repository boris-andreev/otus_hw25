package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HomeworkItem struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Description string             `json:"description" binding:"required" bson:"description"`
}

func (h *HomeworkItem) GetId() primitive.ObjectID {
	return h.Id
}

func (h *HomeworkItem) SetId(id primitive.ObjectID) {
	h.Id = id
}
