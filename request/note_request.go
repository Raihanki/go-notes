package request

type NoteRequest struct {
	Title   string `json:"title" form:"title" validate:"required,min=3,max=150"`
	Content string `json:"content" form:"content" validate:"required,min=3"`
	TopicID uint   `json:"topic_id" form:"topic_id" validate:"required,number"`
}
