package model

type StudyItem struct {
	Id    string `json:"id" bson:"-"`
	Topic string `json:"topic" binding:"required" bson:"topic"`
}

func (s *StudyItem) GetId() string {
	return s.Id
}

func (s *StudyItem) SetId(id string) {
	s.Id = id
}
