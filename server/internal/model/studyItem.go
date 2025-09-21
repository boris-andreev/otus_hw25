package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudyItem struct {
	Id    primitive.ObjectID `json:"id"`
	Topic string             `json:"topic" binding:"required"`
}

func (s *StudyItem) GetId() primitive.ObjectID {
	return s.Id
}

func (s *StudyItem) SetId(id primitive.ObjectID) {
	s.Id = id
}
