// File: internal/service/auth_service.go
package service

import (
    "context"
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "shifa/internal/api/dto"
    "shifa/internal/models"
    "shifa/internal/repository"
)

type AuthService struct {
    userRepo    repository.UserRepository
    jwtSecret   []byte
    tokenExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{
        userRepo:    userRepo,
        jwtSecret:   []byte(jwtSecret),
        tokenExpiry: 24 * time.Hour,
    }
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
    // Check if user exists
    existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, errors.New("user already exists")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Create user
    user := &models.User{
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        Name:         req.Name,
        Role:         req.Role,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // Generate tokens
    return s.generateAuthResponse(user)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
    user, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid credentials")
    }

    return s.generateAuthResponse(user)
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
        "exp":     time.Now().Add(s.tokenExpiry).Unix(),
    })

    tokenString, err := token.SignedString(s.jwtSecret)
    if err != nil {
        return nil, err
    }

    // Generate refresh token (simplified for example)
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
    })

    refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
    if err != nil {
        return nil, err
    }

    return &dto.AuthResponse{
        Token:        tokenString,
        RefreshToken: refreshTokenString,
        User: dto.UserDTO{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Name,
            Role:  user.Role,
        },
    }, nil
}