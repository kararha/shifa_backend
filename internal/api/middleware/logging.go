// File: api/middleware/logging.go
package middleware

import (
    "context"
    "log"
    "net/http"
    "shifa/internal/models"
    "shifa/internal/service"
    "time"
)

type SystemLogMiddleware struct {
    logService *service.SystemLogService
}

func NewSystemLogMiddleware(logService *service.SystemLogService) *SystemLogMiddleware {
    return &SystemLogMiddleware{
        logService: logService,
    }
}

func (m *SystemLogMiddleware) LogSystemAction(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Safely extract user information from context
        userID, ok := r.Context().Value("userID").(int)
        if !ok {
            log.Println("Warning: UserID not found or not an integer")
            userID = 0 // Default userID if not found
        }

        userType, ok := r.Context().Value("userType").(string)
        if !ok {
            log.Println("Warning: UserType not found or not a string")
            userType = "unknown" // Default userType if not found
        }

        // Create system log entry
        logEntry := &models.SystemLog{
            UserID:             userID,
            UserType:           userType,
            ActionType:         r.Method,
            ActionDescription:  r.URL.Path,
            Timestamp:          time.Now(),
            IPAddress:          r.RemoteAddr,
            UserAgent:          r.UserAgent(),
            AdditionalInfo:     make(map[string]interface{}),
        }

        // Create a custom response writer to capture the status code
        rw := &responseWriter{
            ResponseWriter: w,
            statusCode:    http.StatusOK,
        }

        // Defer logging to ensure it happens even if the next handler panics
        defer func() {
            // Add the response status to the log
            logEntry.AdditionalInfo["statusCode"] = rw.statusCode

            // Log the action asynchronously
            go func() {
                if err := m.logService.LogAction(context.Background(), logEntry); err != nil {
                    log.Printf("Error logging system action: %v", err) // Use log.Printf correctly
                }
            }()
        }()

        // Call the next handler
        next.ServeHTTP(rw, r)
    })
}

// Custom response writer to capture status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
    written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
    if !rw.written {
        rw.statusCode = code
        rw.ResponseWriter.WriteHeader(code)
        rw.written = true
    }
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    rw.written = true
    return rw.ResponseWriter.Write(b)
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received %s request for: %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}