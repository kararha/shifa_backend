package service

import (
    "context"
    "fmt"
    "errors"
    "regexp"
    "shifa/internal/models"
    "shifa/internal/repository"
    "strings"
    "unicode"
    "golang.org/x/crypto/bcrypt"
    "time"
)

type UserService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
    return &UserService{
        userRepo: userRepo,
    }
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
    ctx := context.Background()
    
    if !isValidEmail(user.Email) {
        return nil, errors.New("invalid email address")
    }
    if !isValidPassword(user.Password) {
        return nil, errors.New("invalid password")
    }

    hashedPassword, err := hashPassword(user.Password)
    if err != nil {
        return nil, err
    }
    user.PasswordHash = hashedPassword
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    err = s.userRepo.Create(ctx, user)
    if err != nil {
        return nil, err
    }

    return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id int) (*models.User, error) {
    ctx := context.Background()
    return s.userRepo.GetByID(ctx, id)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *models.User) (*models.User, error) {
    ctx := context.Background()
    
    // Verify user exists
    existing, err := s.userRepo.GetByID(ctx, user.ID)
    if err != nil {
        return nil, err
    }

    // If password is provided, validate and hash it
    if user.Password != "" {
        if !isValidPassword(user.Password) {
            return nil, errors.New("invalid password")
        }
        hashedPassword, err := hashPassword(user.Password)
        if err != nil {
            return nil, err
        }
        user.PasswordHash = hashedPassword
    } else {
        // Keep existing password hash if no new password provided
        user.PasswordHash = existing.PasswordHash
    }

    user.UpdatedAt = time.Now()
    
    err = s.userRepo.Update(ctx, user)
    if err != nil {
        return nil, err
    }

    return user, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id int) error {
    ctx := context.Background()
    return s.userRepo.Delete(ctx, id)
}

// Register handles the registration of a new user
func (s *UserService) Register(ctx context.Context, user models.User) error {
    if !isValidEmail(user.Email) {
        return errors.New("invalid email address")
    }
    if !isValidPassword(user.Password) {
        return errors.New("password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character")
    }
    
    hashedPassword, err := hashPassword(user.Password)
    if err != nil {
        return errors.New("failed to hash password")
    }
    user.PasswordHash = hashedPassword
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    return s.userRepo.Create(ctx, &user)
}

// Authenticate handles the login process
func (s *UserService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    
    if !verifyPassword(user.PasswordHash, password) {
        return nil, errors.New("invalid credentials")
    }
    return user, nil
}

// List retrieves a paginated list of users
func (s *UserService) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
    return s.userRepo.List(ctx, offset, limit)
}

// Helper functions remain the same
func isValidPassword(password string) bool {
    var (
        hasUpperCase, hasLowerCase, hasNumber, hasSpecialChar bool
        specialChars = "!@#$%^&*(),.?\":{}|<>" // Your current special chars
    )
    
    // Add debugging
    // fmt.Printf("Checking password: %s\n", password)
    
    if len(password) < 8 {
        fmt.Println("Password length less than 8")
        return false
    }
    
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpperCase = true
            // fmt.Printf("Found uppercase: %c\n", char)
        case unicode.IsLower(char):
            hasLowerCase = true
            // fmt.Printf("Found lowercase: %c\n", char)
        case unicode.IsNumber(char):
            hasNumber = true
            // fmt.Printf("Found number: %c\n", char)
        case strings.ContainsRune(specialChars, char):
            hasSpecialChar = true
            // fmt.Printf("Found special char: %c\n", char)
        }
    }
    
    fmt.Printf("Validation results:\nUppercase: %v\nLowercase: %v\nNumber: %v\nSpecial: %v\n",
        hasUpperCase, hasLowerCase, hasNumber, hasSpecialChar)
    
    return hasUpperCase && hasLowerCase && hasNumber && hasSpecialChar
}

func hashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func verifyPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func isValidEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}