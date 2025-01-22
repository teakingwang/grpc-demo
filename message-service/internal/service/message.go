package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/teakingwang/grpc-demo/pkg/errors"
	pb "github.com/teakingwang/grpc-demo/proto/message"
)

type MessageService struct {
	pb.UnimplementedMessageServiceServer
	db *sql.DB
}

func NewMessageService(db *sql.DB) *MessageService {
	return &MessageService{db: db}
}

func (s *MessageService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if req.FromUserId <= 0 || req.ToUserId <= 0 {
		return nil, errors.ErrInvalidInput
	}

	result, err := s.db.ExecContext(ctx,
		"INSERT INTO messages (from_user_id, to_user_id, content, created_at) VALUES (?, ?, ?, ?)",
		req.FromUserId, req.ToUserId, req.Content, time.Now().Unix())
	if err != nil {
		return nil, errors.FromError(err)
	}

	id, _ := result.LastInsertId()
	return &pb.SendMessageResponse{
		MessageId: id,
	}, nil
}

func (s *MessageService) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, from_user_id, to_user_id, content, created_at 
         FROM messages 
         WHERE to_user_id = ? 
         ORDER BY created_at DESC`,
		req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*pb.Message
	for rows.Next() {
		msg := &pb.Message{}
		err := rows.Scan(&msg.Id, &msg.FromUserId, &msg.ToUserId, &msg.Content, &msg.CreatedAt)
		if err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	return &pb.GetMessagesResponse{
		Messages: messages,
	}, nil
}
