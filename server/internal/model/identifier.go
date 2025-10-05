package model

type Identifier interface {
	GetId() string
	SetId(string)
}

type ItemWithId interface {
	*HomeworkItem | *StudyItem | *WorkoutItem
	Identifier
}
