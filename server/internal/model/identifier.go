package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Identifier interface {
	GetId() primitive.ObjectID
	SetId(primitive.ObjectID)
}

type ItemWithId interface {
	*HomeworkItem | *StudyItem | *WorkoutItem
	Identifier
}
