package model

type HomeworkItem struct {
	Id          string `json:"id" bson:"-"`
	Description string `json:"description" binding:"required" bson:"description"`
}

func (h *HomeworkItem) GetId() string {
	return h.Id
}

func (h *HomeworkItem) SetId(id string) {
	h.Id = id
}
