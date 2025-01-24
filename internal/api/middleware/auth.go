// File: internal/api/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
    jwtSecret []byte
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
    return &AuthMiddleware{
        jwtSecret: []byte(jwtSecret),
    }
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return m.jwtSecret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        // Add claims to context
        ctx := context.WithValue(r.Context(), "userID", int(claims["user_id"].(float64)))
        ctx = context.WithValue(ctx, "userRole", claims["role"].(string))

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}