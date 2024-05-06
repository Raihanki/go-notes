package resources

import "time"

type NoteResource struct {
	ID        uint               `json:"id"`
	User      UserPublicResource `json:"user"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	Topic     TopicResource      `json:"topic"`
	CreatedAt time.Time          `json:"created_at"`
}
