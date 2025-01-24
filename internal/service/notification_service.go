// service/notification_service.go
package service

import (
    "context"
    "errors"
    "time"
    
    "github.com/sirupsen/logrus" // Import logrus
    "shifa/internal/models"       // Import models
    "shifa/internal/repository"   // Import repository
)

type NotificationService interface {
    CreateNotification(ctx context.Context, notification *models.Notification) error
    GetNotificationsByUserID(ctx context.Context, userID int, page, pageSize int) ([]*models.Notification, error)
    MarkNotificationAsRead(ctx context.Context, notificationID int) error
    GetUnreadNotificationCount(ctx context.Context, userID int) (int, error)
    SendAppointmentReminder(ctx context.Context, appointmentID int) error
}

type notificationService struct {
    notificationRepo repository.NotificationRepository
    logger           *logrus.Logger // Use *logrus.Logger
}

// Constructor for NotificationService
func NewNotificationService(notificationRepo repository.NotificationRepository, logger *logrus.Logger) NotificationService {
    return &notificationService{
        notificationRepo: notificationRepo,
        logger:           logger,
    }
}

// CreateNotification creates a new notification
func (s *notificationService) CreateNotification(ctx context.Context, notification *models.Notification) error {
    notification.CreatedAt = time.Now()
    notification.IsRead = false

    err := s.notificationRepo.Create(ctx, notification)
    if err != nil {
        s.logger.WithFields(logrus.Fields{
            "error":        err,
            "notification": notification,
        }).Error("Failed to create notification")
        return errors.New("failed to create notification")
    }

    return nil
}

// GetNotificationsByUserID retrieves notifications for a user
func (s *notificationService) GetNotificationsByUserID(ctx context.Context, userID int, page, pageSize int) ([]*models.Notification, error) {
    offset := (page - 1) * pageSize
    notifications, err := s.notificationRepo.GetByUserID(ctx, userID, pageSize, offset)
    if err != nil {
        s.logger.WithFields(logrus.Fields{
            "error":   err,
            "userID":  userID,
        }).Error("Failed to get notifications by user ID")
        return nil, errors.New("failed to fetch notifications")
    }
    return notifications, nil
}

// MarkNotificationAsRead marks a notification as read
func (s *notificationService) MarkNotificationAsRead(ctx context.Context, notificationID int) error {
    err := s.notificationRepo.MarkAsRead(ctx, notificationID)
    if err != nil {
        s.logger.WithFields(logrus.Fields{
            "error":          err,
            "notificationID": notificationID,
        }).Error("Failed to mark notification as read")
        return errors.New("failed to update notification status")
    }
    return nil
}

// GetUnreadNotificationCount retrieves the unread notification count for a user
func (s *notificationService) GetUnreadNotificationCount(ctx context.Context, userID int) (int, error) {
    count, err := s.notificationRepo.GetUnreadCount(ctx, userID)
    if err != nil {
        s.logger.WithFields(logrus.Fields{
            "error":   err,
            "userID":  userID,
        }).Error("Failed to get unread notification count")
        return 0, errors.New("failed to fetch unread notification count")
    }
    return count, nil
}

// SendAppointmentReminder sends an appointment reminder notification
func (s *notificationService) SendAppointmentReminder(ctx context.Context, appointmentID int) error {
    // Logic to fetch appointment details and create a reminder notification
    notification := &models.Notification{ // Use models.Notification
        UserID:           0, // Set the actual user ID
        NotificationType: "appointment_reminder",
        Message:          "You have an upcoming appointment.",
    }
    return s.CreateNotification(ctx, notification)
}