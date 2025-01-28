package service

import (
	"context"
	"errors"
	"fmt"

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
	if _, err := a.authRepo.ReadUserByPhone(ctx, creds.Phonenumber); err == nil {
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

func (a *authService) ReadUserByPhone(ctx context.Context, phonenumber string) (*domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (a *authService) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (a *authService) DeleteUser(ctx context.Context, id string) error {
	panic("not implemented") // TODO: Implement
}

func (a *authService) CreateSession(ctx context.Context, userid string) (*domain.Session, error) {
	panic("not implemented") // TODO: Implement
}

func (a *authService) ReadSession(ctx context.Context, token string) (*domain.Session, error) {
	panic("not implemented") // TODO: Implement
}

func (a *authService) DeleteSession(ctx context.Context, token string) error {
	panic("not implemented") // TODO: Implement
}
func NewAuthService(ar port.AuthRepository) port.AuthService {
	return &authService{authRepo: ar}
}
