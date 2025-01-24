// // repository/mysql/notification_repo.go
// package mysql

// import (
//     "context"
//     "database/sql"
//     "errors"
//     // "time"

//     "shifa/internal/models"
//     "shifa/pkg/logger"
// )

// // DbInterface is an interface that both *sql.DB and *sql.Tx satisfy
// type DbInterface interface {
//     ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
//     QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
//     QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
// }

// // NotificationRepo represents the MySQL repository for notification-related database operations
// type NotificationRepo struct {
//     db     DbInterface
//     logger logger.Logger
// }

// // NewNotificationRepo creates a new NotificationRepo instance
// func NewNotificationRepo(db DbInterface, log logger.Logger) *NotificationRepo {
//     return &NotificationRepo{
//         db:     db,
//         logger: log,
//     }
// }

// func (r *NotificationRepo) Create(ctx context.Context, notification *models.Notification) error {
//     query := `INSERT INTO notifications (user_id, notification_type, message, is_read, created_at)
//               VALUES (?, ?, ?, ?, ?)`
    
//     _, err := r.db.ExecContext(ctx, query, notification.UserID, notification.NotificationType,
//         notification.Message, notification.IsRead, notification.CreatedAt)
//     if err != nil {
//         r.logger.Error("Failed to create notification", "error", err)
//         return errors.New("failed to create notification")
//     }
//     return nil
// }

// func (r *NotificationRepo) GetByUserID(ctx context.Context, userID int, limit, offset int) ([]*models.Notification, error) {
//     query := `SELECT id, user_id, notification_type, message, is_read, created_at
//               FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
    
//     rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
//     if err != nil {
//         r.logger.Error("Failed to get notifications", "error", err, "userID", userID)
//         return nil, errors.New("failed to get notifications")
//     }
//     defer rows.Close()
    
//     var notifications []*models.Notification
//     for rows.Next() {
//         var notification models.Notification
//         err := rows.Scan(
//             &notification.ID, &notification.UserID, &notification.NotificationType,
//             &notification.Message, &notification.IsRead, &notification.CreatedAt,
//         )
//         if err != nil {
//             r.logger.Error("Failed to scan notification row", "error", err)
//             return nil, errors.New("failed to process notification data")
//         }
//         notifications = append(notifications, &notification)
//     }
    
//     if err = rows.Err(); err != nil {
//         r.logger.Error("Error after scanning notifications", "error", err)
//         return nil, errors.New("error processing notifications")
//     }
    
//     return notifications, nil
// }

// func (r *NotificationRepo) MarkAsRead(ctx context.Context, notificationID int) error {
//     query := `UPDATE notifications SET is_read = true WHERE id = ?`
    
//     result, err := r.db.ExecContext(ctx, query, notificationID)
//     if err != nil {
//         r.logger.Error("Failed to mark notification as read", "error", err, "notificationID", notificationID)
//         return errors.New("failed to update notification")
//     }
    
//     rowsAffected, err := result.RowsAffected()
//     if err != nil {
//         r.logger.Error("Failed to get rows affected", "error", err)
//         return errors.New("failed to confirm notification update")
//     }
    
//     if rowsAffected == 0 {
//         return errors.New("notification not found")
//     }
    
//     return nil
// }

// func (r *NotificationRepo) GetUnreadCount(ctx context.Context, userID int) (int, error) {
//     query := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = false`
    
//     var count int
//     err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
//     if err != nil {
//         r.logger.Error("Failed to get unread notification count", "error", err, "userID", userID)
//         return 0, errors.New("failed to get unread notification count")
//     }
//     return count, nil
// }

// // WithTx returns a new NotificationRepo with the given transaction
// func (r *NotificationRepo) WithTx(tx *sql.Tx) *NotificationRepo {
//     return &NotificationRepo{
//         db:     tx,
//         logger: r.logger,
//     }
// }



// repository/mysql/notification_repo.go
package mysql

import (
    "context"
    "database/sql"
    "errors"

    "github.com/sirupsen/logrus"
    "shifa/internal/models"
)

// DbInterface is an interface that both *sql.DB and *sql.Tx satisfy
type DbInterface interface {
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// NotificationRepo represents the MySQL repository for notification-related database operations
type NotificationRepo struct {
    db     DbInterface
    logger *logrus.Logger
}

// NewNotificationRepo creates a new NotificationRepo instance
func NewNotificationRepo(db DbInterface, log *logrus.Logger) *NotificationRepo {
    return &NotificationRepo{
        db:     db,
        logger: log,
    }
}

func (r *NotificationRepo) Create(ctx context.Context, notification *models.Notification) error {
    query := `INSERT INTO notifications (user_id, notification_type, message, is_read, created_at)
              VALUES (?, ?, ?, ?, ?)`
    
    _, err := r.db.ExecContext(ctx, query, notification.UserID, notification.NotificationType,
        notification.Message, notification.IsRead, notification.CreatedAt)
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "error":              err,
            "user_id":           notification.UserID,
            "notification_type":  notification.NotificationType,
        }).Error("Failed to create notification")
        return errors.New("failed to create notification")
    }
    return nil
}

func (r *NotificationRepo) GetByUserID(ctx context.Context, userID int, limit, offset int) ([]*models.Notification, error) {
    query := `SELECT id, user_id, notification_type, message, is_read, created_at
              FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
    
    rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "error": err,
            "user_id": userID,
        }).Error("Failed to get notifications")
        return nil, errors.New("failed to get notifications")
    }
    defer rows.Close()
    
    var notifications []*models.Notification
    for rows.Next() {
        var notification models.Notification
        err := rows.Scan(
            &notification.ID, &notification.UserID, &notification.NotificationType,
            &notification.Message, &notification.IsRead, &notification.CreatedAt,
        )
        if err != nil {
            r.logger.Error("Failed to scan notification row", "error", err)
            return nil, errors.New("failed to process notification data")
        }
        notifications = append(notifications, &notification)
    }
    
    if err = rows.Err(); err != nil {
        r.logger.Error("Error after scanning notifications", "error", err)
        return nil, errors.New("error processing notifications")
    }
    
    return notifications, nil
}

func (r *NotificationRepo) MarkAsRead(ctx context.Context, notificationID int) error {
    query := `UPDATE notifications SET is_read = true WHERE id = ?`
    
    result, err := r.db.ExecContext(ctx, query, notificationID)
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "error": err,
            "notificationID": notificationID,
        }).Error("Failed to mark notification as read")
        return errors.New("failed to update notification")
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        r.logger.Error("Failed to get rows affected", "error", err)
        return errors.New("failed to confirm notification update")
    }
    
    if rowsAffected == 0 {
        return errors.New("notification not found")
    }
    
    return nil
}

func (r *NotificationRepo) GetUnreadCount(ctx context.Context, userID int) (int, error) {
    query := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = false`
    
    var count int
    err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "error": err,
            "userID": userID,
        }).Error("Failed to get unread notification count")
        return 0, errors.New("failed to get unread notification count")
    }
    return count, nil
}

// WithTx returns a new NotificationRepo with the given transaction
func (r *NotificationRepo) WithTx(tx *sql.Tx) *NotificationRepo {
    return &NotificationRepo{
        db:     tx,
        logger: r.logger,
    }
}