package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "shifa/internal/models"
    "shifa/internal/repository"
    "shifa/internal/service"
)

type ChatMessageHandler struct {
    chatService service.ChatService  // Use the interface type directly
}

func NewChatMessageHandler(chatService service.ChatService) *ChatMessageHandler {
    return &ChatMessageHandler{chatService: chatService}
}

func (h *ChatMessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
    var msg models.ChatMessage
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    ctx := r.Context()
    // Convert models.ChatMessage to repository.ChatMessage
    repoMessage := repository.ChatMessage{
        ID:             msg.ID,
        ConsultationID: msg.ConsultationID,
        SenderType:     msg.SenderType,
        SenderID:       msg.SenderID,
        Message:        msg.Message,
        SentAt:         msg.SentAt,
        IsRead:         msg.IsRead,
    }
    if err := h.chatService.SendMessage(ctx, &repoMessage); err != nil {
        http.Error(w, "Failed to send message", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}

func (h *ChatMessageHandler) GetMessagesByConsultation(w http.ResponseWriter, r *http.Request) {
    consultationID, err := strconv.Atoi(r.URL.Query().Get("consultation_id"))
    if err != nil {
        http.Error(w, "Invalid consultation ID", http.StatusBadRequest)
        return
    }
    ctx := r.Context()
    messages, err := h.chatService.GetMessagesByConsultationID(ctx, consultationID, 1, 10)
    if err != nil {
        http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(messages)
}


func (h *ChatMessageHandler) MarkMessageAsRead(w http.ResponseWriter, r *http.Request, messageID int) {
    ctx := r.Context()
    if err := h.chatService.MarkMessageAsRead(ctx, messageID); err != nil {
        http.Error(w, "Failed to mark message as read", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (h *ChatMessageHandler) GetUnreadMessageCount(w http.ResponseWriter, r *http.Request, userID int) {
    ctx := r.Context()
    count, err := h.chatService.GetUnreadMessageCount(ctx, userID)
    if err != nil {
        http.Error(w, "Failed to get unread message count", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]int{"unread_count": count})
}
