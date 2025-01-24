package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "shifa/internal/models"
    "shifa/internal/service"
)

type NotificationHandler struct {
    notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
    return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
    var notification models.Notification
    if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Directly use models.Notification since the service expects this type
    if err := h.notificationService.CreateNotification(r.Context(), &notification); err != nil {
        http.Error(w, "Failed to create notification", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
    userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Default page and pageSize values
    page := 1
    pageSize := 10

    // Get page from query parameters if provided
    if pageStr := r.URL.Query().Get("page"); pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }

    // Get pageSize from query parameters if provided
    if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
        if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
            pageSize = ps
        }
    }

    notifications, err := h.notificationService.GetNotificationsByUserID(r.Context(), userID, page, pageSize)
    if err != nil {
        http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
        return
    }

    // Convert model notifications to JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(notifications) // Directly encode notifications
}

func (h *NotificationHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    notificationID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid notification ID", http.StatusBadRequest)
        return
    }

    if err := h.notificationService.MarkNotificationAsRead(r.Context(), notificationID); err != nil {
        http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
    userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    count, err := h.notificationService.GetUnreadNotificationCount(r.Context(), userID)
    if err != nil {
        http.Error(w, "Failed to get unread count", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"count": count})
}

func (h *NotificationHandler) SendAppointmentReminder(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    appointmentID, err := strconv.Atoi(vars["appointmentId"])
    if err != nil {
        http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
        return
    }

    if err := h.notificationService.SendAppointmentReminder(r.Context(), appointmentID); err != nil {
        http.Error(w, "Failed to send appointment reminder", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}