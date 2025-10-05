package model

type WorkoutItem struct {
	Id     string `json:"id"  bson:"-"`
	Target string `json:"target" binding:"required" bson:"target"`
}

func (w *WorkoutItem) GetId() string {
	return w.Id
}

func (w *WorkoutItem) SetId(id string) {
	w.Id = id
}
