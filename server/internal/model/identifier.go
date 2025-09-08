package model

type Identifier interface {
	GetId() int
	SetId(int)
}

type ItemWithId interface {
	*HomeworkItem | *StudyItem | *WorkoutItem
	Identifier
}
