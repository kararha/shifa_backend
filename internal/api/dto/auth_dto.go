// File: internal/api/dto/auth_dto.go
package dto

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
    Role     string `json:"role" validate:"required,oneof=doctor patient home_care_provider admin"`
}

type AuthResponse struct {
    Token        string `json:"token"`
    RefreshToken string `json:"refresh_token"`
    User         UserDTO `json:"user"`
}

type UserDTO struct {
    ID    int    `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
    Role  string `json:"role"`
}
