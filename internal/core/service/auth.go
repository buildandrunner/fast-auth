package service

import (
	"context"
	"errors"
	"log"

	"github.com/mar-cial/space-auth/internal/core/domain"
	"github.com/mar-cial/space-auth/internal/core/port"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInternalError     = errors.New("internal server error")
	ErrUserNotFound      = errors.New("user not found")
	ErrSessionNotFound   = errors.New("session not found")
)

type authService struct {
	authRepo port.AuthRepository
}

func (a *authService) Register(ctx context.Context, phonenumber string) (*domain.User, error) {
	// Check if the user already exists
	existingUser, err := a.authRepo.FindUserByPhone(ctx, phonenumber)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Proceed with registration
	newUser := &domain.User{
		ID:          generateUniqueID(),
		Phonenumber: phonenumber,
	}

	if err := a.authRepo.SaveUser(ctx, newUser); err != nil {
		return nil, ErrInternalError
	}

	session, err := createSession(newUser.ID, 124)
	if err != nil {
		return nil, ErrInternalError
	}

	if err := a.authRepo.CreateSession(ctx, session); err != nil {
		return nil, ErrInternalError
	}

	return newUser, nil
}

func (a *authService) Login(ctx context.Context, phonenumber string) (*domain.User, error) {
	// Verify that the user exists
	user, err := a.authRepo.FindUserByPhone(ctx, phonenumber)
	if err != nil {
		log.Println(err)
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalError
	}

	// Create a session for the user
	session, err := createSession(user.ID, 124)
	if err != nil {
		log.Println(err)
		return nil, ErrInternalError
	}

	if err := a.authRepo.CreateSession(ctx, session); err != nil {
		log.Println(err)
		return nil, ErrInternalError
	}

	return user, nil
}

func (a *authService) Logout(ctx context.Context, token string) error {
	err := a.authRepo.DeleteSession(ctx, token)
	if err != nil {
		log.Println(err)
		if errors.Is(err, ErrSessionNotFound) {
			return ErrSessionNotFound
		}
		return ErrInternalError
	}

	return nil
}

// NewAuthService creates and returns an instance of authService.
func NewAuthService(authRepo port.AuthRepository) port.AuthService {
	return &authService{authRepo: authRepo}
}
