package request

type NoteRequest struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	TopicID uint   `json:"topic_id" form:"topic_id"`
}
