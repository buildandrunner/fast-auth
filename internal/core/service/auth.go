package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mar-cial/space-auth/internal/core/domain"
	"github.com/mar-cial/space-auth/internal/core/port"
)

type authService struct {
	authRepo port.AuthRepository
}

var (
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
)

// CreateUser with Argon2id password hashing
func (a *authService) CreateUser(ctx context.Context, creds domain.Credentials) (*domain.User, error) {
	// Check for existing user
	foundUser, err := a.authRepo.ReadUserByPhone(ctx, creds.Phonenumber)
	if err != nil {
		return nil, err
	}

	if foundUser != nil {
		return nil, ErrUserExists
	}

	// Generate password hash
	encodedHash, err := generateFromPassword(creds.Password, defaultArgon2Params())
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	// Create domain user
	user := &domain.User{
		ID:          generateUniqueID(),
		Phonenumber: creds.Phonenumber,
		Password:    encodedHash,
	}

	// Save to repository
	if _, err := a.authRepo.SaveUser(ctx, *user); err != nil {
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	return user, nil
}

// ValidateUser credentials with Argon2id
func (a *authService) ValidateUser(ctx context.Context, creds domain.Credentials) (bool, error) {
	user, err := a.authRepo.ReadUserByPhone(ctx, creds.Phonenumber)
	if err != nil {
		if errors.Is(err, port.ErrUserNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("validation failed: %w", err)
	}

	match, err := comparePasswordAndHash(creds.Password, user.Password)
	if err != nil {
		return false, fmt.Errorf("password comparison failed: %w", err)
	}

	return match, nil
}

// ReadUserById fetches user by ID
func (a *authService) ReadUserById(ctx context.Context, id string) (*domain.User, error) {
	user, err := a.authRepo.ReadUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, port.ErrUserNotFound) {
			return nil, port.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	return user, nil
}

// service/auth_service.go
func (a *authService) ReadUserByPhone(ctx context.Context, phonenumber string) (*domain.User, error) {
	user, err := a.authRepo.ReadUserByPhone(ctx, phonenumber)
	if err != nil {
		if errors.Is(err, port.ErrUserNotFound) {
			return nil, port.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to read user by phone: %w", err)
	}
	return user, nil
}

func (a *authService) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	// Verify existing user
	existingUser, err := a.authRepo.ReadUserByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("user verification failed: %w", err)
	}

	// Check phone number availability if changing
	if user.Phonenumber != existingUser.Phonenumber {
		if _, err := a.authRepo.ReadUserByPhone(ctx, user.Phonenumber); err == nil {
			return nil, ErrUserExists
		}
	}

	// Perform atomic update
	updatedUser, err := a.authRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("update operation failed: %w", err)
	}

	return updatedUser, nil
}

func (a *authService) DeleteUser(ctx context.Context, id string) error {
	// Get user details first
	user, err := a.authRepo.ReadUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user lookup failed: %w", err)
	}

	// Perform atomic deletion
	if err := a.authRepo.DeleteUser(ctx, *user); err != nil {
		return fmt.Errorf("deletion failed: %w", err)
	}

	return nil
}

func (a *authService) CreateSession(ctx context.Context, userid string) (*domain.Session, error) {
	// Generate new session with 24h duration
	session, err := generateSession(userid, 24)
	if err != nil {
		return nil, fmt.Errorf("session generation failed: %w", err)
	}

	// Save to repository
	if _, err := a.authRepo.SaveSession(ctx, *session, userid); err != nil {
		return nil, fmt.Errorf("session persistence failed: %w", err)
	}

	return session, nil
}

func (a *authService) ReadSession(ctx context.Context, token string) (*domain.Session, error) {
	session, err := a.authRepo.FindSessionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, port.ErrUserNotFound) {
			return nil, port.ErrUserNotFound
		}
		return nil, fmt.Errorf("session retrieval failed: %w", err)
	}

	// Validate session expiration
	if time.Now().After(session.ExpiresAt) {
		// Auto-cleanup expired session
		_ = a.authRepo.DeleteSession(ctx, token)
		return nil, port.ErrSessionExpired
	}

	return session, nil
}

func (a *authService) DeleteSession(ctx context.Context, token string) error {
	if err := a.authRepo.DeleteSession(ctx, token); err != nil {
		return fmt.Errorf("session deletion failed: %w", err)
	}
	return nil
}

func NewAuthService(ar port.AuthRepository) port.AuthService {
	return &authService{authRepo: ar}
}
