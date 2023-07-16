package core

import "time"

type Message struct {
	Id          int64     `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	FromUserId  int64     `json:"user_id" db:"user_id"`
	ToRoomId    int64     `json:"room_id" db:"room_id"`
	Content     string    `json:"content" db:"content"`
	ContentType int8      `json:"content_type" db:"content_type"`
	Url         string    `json:"url" db:"url"`
}
