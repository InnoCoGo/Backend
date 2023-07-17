package core

type Message struct {
	Id           int64  `json:"id" db:"id"`
	FromUsername string `json:"fromUsername" db:"username"`
	FromUserId   int64  `json:"fromUserId" db:"user_id"`
	ToRoomId     int64  `json:"toRoomId" db:"room_id"`
	Content      string `json:"content" db:"content"`
	ContentType  int8   `json:"contentType" db:"content_type"`
	Url          string `json:"url" db:"url"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}

const (
	TEXT = iota
	IMAGE
	VIDEO
	AUDIO
	LOCATION
	DOCUMENT
	CONTACT
)
