package repository

import (
	"fmt"
	"time"

	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/protocol"
	"github.com/jmoiron/sqlx"
)

type MessagePostgres struct {
	db *sqlx.DB
}

func NewMessagePostgres(db *sqlx.DB) *MessagePostgres {
	return &MessagePostgres{db: db}
}

func (r *MessagePostgres) Save(message protocol.Message) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s 
	(created_at, user_id, room_id, content, content_type, url)
						VAlUES
	($1, $2, $3, $4, $5, $6) RETURNING id`, messagesTable)

	msg := core.Message{
		CreatedAt:   time.Now(),
		FromUserId:  message.FromId,
		ToRoomId:    message.ToRoomId,
		Content:     message.Content,
		ContentType: int8(message.ContentType),
		Url:         message.Url,
	}

	var id int64
	row := r.db.QueryRow(query, msg.CreatedAt, msg.FromUserId, msg.ToRoomId, msg.Content, msg.ContentType, msg.Url)
	err := row.Scan(&id)
	return id, err
}
