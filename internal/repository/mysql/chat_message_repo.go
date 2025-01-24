// repository/mysql/chat_message_repo.go

package mysql

import (
	"context"
	"database/sql"

	"shifa/internal/repository"
)

type mysqlChatMessageRepo struct {
	db *sql.DB
}

func NewMySQLChatMessageRepo(db *sql.DB) repository.ChatMessageRepository {
	return &mysqlChatMessageRepo{db: db}
}

func (r *mysqlChatMessageRepo) Create(ctx context.Context, message *repository.ChatMessage) error {
	query := `INSERT INTO chat_messages (consultation_id, sender_type, sender_id, message, sent_at, is_read)
			  VALUES (?, ?, ?, ?, ?, ?)`
	
	_, err := r.db.ExecContext(ctx, query, message.ConsultationID, message.SenderType, message.SenderID,
		message.Message, message.SentAt, message.IsRead)
	return err
}

func (r *mysqlChatMessageRepo) GetByConsultationID(ctx context.Context, consultationID int, limit, offset int) ([]*repository.ChatMessage, error) {
	query := `SELECT id, consultation_id, sender_type, sender_id, message, sent_at, is_read
			  FROM chat_messages WHERE consultation_id = ? ORDER BY sent_at DESC LIMIT ? OFFSET ?`
	
	rows, err := r.db.QueryContext(ctx, query, consultationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []*repository.ChatMessage
	for rows.Next() {
		var message repository.ChatMessage
		err := rows.Scan(
			&message.ID, &message.ConsultationID, &message.SenderType, &message.SenderID,
			&message.Message, &message.SentAt, &message.IsRead,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	
	return messages, nil
}

func (r *mysqlChatMessageRepo) MarkAsRead(ctx context.Context, messageID int) error {
	query := `UPDATE chat_messages SET is_read = true WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, messageID)
	return err
}

func (r *mysqlChatMessageRepo) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	query := `SELECT COUNT(*) FROM chat_messages WHERE sender_id != ? AND is_read = false`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}