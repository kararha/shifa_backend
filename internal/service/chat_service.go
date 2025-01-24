package service

import (
    "context"
    "errors"
    "time"
    "github.com/sirupsen/logrus"
    "shifa/internal/repository"
)

// ChatService interface defines the contract for chat operations
type ChatService interface {
    SendMessage(ctx context.Context, message *repository.ChatMessage) error
    GetMessagesByConsultationID(ctx context.Context, consultationID int, page, pageSize int) ([]*repository.ChatMessage, error)
    MarkMessageAsRead(ctx context.Context, messageID int) error
    GetUnreadMessageCount(ctx context.Context, userID int) (int, error)
}

type chatService struct {
    chatRepo repository.ChatMessageRepository
    log      *logrus.Entry // Change this to *logrus.Entry
}

// NewChatService constructor
func NewChatService(
    chatRepo repository.ChatMessageRepository, 
    log *logrus.Logger, // Accept *logrus.Logger
) ChatService {
    // Create a service-specific logger with additional context
    serviceLogger := log.WithField("service", "chat_service")
    
    return &chatService{
        chatRepo: chatRepo,
        log:      serviceLogger, // Use the *logrus.Entry
    }
}

func (s *chatService) SendMessage(ctx context.Context, message *repository.ChatMessage) error {
    methodLogger := s.log.WithFields(logrus.Fields{
        "method":         "SendMessage",
        "senderID":      message.SenderID,
        "senderType":    message.SenderType,
        "consultationID": message.ConsultationID,
    })

    // Set message metadata
    message.SentAt = time.Now()
    message.IsRead = false

    // Attempt to create message
    err := s.chatRepo.Create(ctx, message)
    if err != nil {
        methodLogger.Error("Failed to send message", "error", err.Error(), "details", "message creation failed")
        return errors.New("failed to send message")
    }

    methodLogger.Info("Message sent successfully", "messageID", message.ID)
    return nil
}

func (s *chatService) GetMessagesByConsultationID(ctx context.Context, consultationID int, page, pageSize int) ([]*repository.ChatMessage, error) {
    methodLogger := s.log.WithFields(logrus.Fields{
        "method":          "GetMessagesByConsultationID",
        "consultationID":  consultationID,
        "page":            page,
        "pageSize":        pageSize,
    })

    // Calculate offset
    offset := (page - 1) * pageSize

    // Fetch messages
    messages, err := s.chatRepo.GetByConsultationID(ctx, consultationID, pageSize, offset)
    if err != nil {
        methodLogger.Error("Failed to retrieve messages", "error", err.Error(), "details", "database query failed")
        return nil, errors.New("failed to fetch messages")
    }

    methodLogger.Debug("Messages retrieved successfully", "messageCount", len(messages))
    return messages, nil
}

func (s *chatService) MarkMessageAsRead(ctx context.Context, messageID int) error {
    methodLogger := s.log.WithFields(logrus.Fields{
        "method":     "MarkMessageAsRead",
        "messageID":  messageID,
    })

    // Attempt to mark message as read
    err := s.chatRepo.MarkAsRead(ctx, messageID)
    if err != nil {
        methodLogger.Error("Failed to mark message as read", "error", err.Error(), "details", "update operation failed")
        return errors.New("failed to update message status")
    }

    methodLogger.Info("Message marked as read successfully")
    return nil
}

func (s *chatService) GetUnreadMessageCount(ctx context.Context, userID int) (int, error) {
    methodLogger := s.log.WithFields(logrus.Fields{
        "method":  "GetUnreadMessageCount",
        "userID":  userID,
    })

    // Fetch unread message count
    count, err := s.chatRepo.GetUnreadCount(ctx, userID)
    if err != nil {
        methodLogger.Error("Failed to get unread message count", "error", err.Error(), "details", "count retrieval failed")
        return 0, errors.New("failed to fetch unread message count")
    }

    methodLogger.Debug("Unread message count retrieved", "unreadCount", count)
    return count, nil
}