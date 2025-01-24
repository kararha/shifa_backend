// middleware/cors.go
package middleware

import (
    "net/http"
    "github.com/rs/cors"
)

func CORSMiddleware() func(http.Handler) http.Handler {
    return cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"}, // Svelte dev servers
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }).Handler
}