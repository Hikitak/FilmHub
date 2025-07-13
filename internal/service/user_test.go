package service

import (
    "context"
    "testing"

    "filmhub/internal/models"
)

// stubUserRepo is an in-memory implementation of repository.UserRepository
// used exclusively in unit tests.
// Only methods needed by AuthService are implemented.
type stubUserRepo struct {
    users map[string]*models.User // keyed by email
}

func newStubUserRepo() *stubUserRepo {
    return &stubUserRepo{users: make(map[string]*models.User)}
}

func (s *stubUserRepo) Create(_ context.Context, user *models.User) error {
    s.users[user.Email] = user
    return nil
}

func (s *stubUserRepo) FindByEmail(_ context.Context, email string) (*models.User, error) {
    if u, ok := s.users[email]; ok {
        return u, nil
    }
    return nil, context.Canceled // use any non-nil error to indicate not found
}

func TestAuthService_Register_And_Login(t *testing.T) {
    repo := newStubUserRepo()
    svc := NewAuthService(repo)

    ctx := context.Background()
    user := &models.User{
        Username: "john",
        Email:    "john@example.com",
        Password: "s3cr3tPwd",
    }

    // Register should hash password and store the user.
    if err := svc.Register(ctx, user); err != nil {
        t.Fatalf("register failed: %v", err)
    }
    stored, _ := repo.FindByEmail(ctx, user.Email)
    if stored == nil {
        t.Fatalf("user not stored after register")
    }
    if stored.Password == "s3cr3tPwd" {
        t.Errorf("password was not hashed")
    }

    // Login with correct credentials should return token.
    token, err := svc.Login(ctx, user.Email, "s3cr3tPwd")
    if err != nil {
        t.Fatalf("login failed: %v", err)
    }
    if token == "" {
        t.Errorf("expected non-empty token")
    }

    // Login with wrong password should error.
    if _, err := svc.Login(ctx, user.Email, "wrong"); err == nil {
        t.Errorf("expected error for wrong password, got nil")
    }
} 