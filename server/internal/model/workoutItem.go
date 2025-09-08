package model

type WorkoutItem struct {
	Id     int    `json:"id"`
	Target string `json:"target" binding:"required"`
}

func (w *WorkoutItem) GetId() int {
	return w.Id
}

func (w *WorkoutItem) SetId(id int) {
	w.Id = id
}
