package port

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mar-cial/space-auth/internal/core/domain"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

// auth core
type AuthHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type AuthService interface {
	UserService
	SessionService
}

type AuthRepository interface {
	UserRepository
	SessionRepository
}

// service layer
type UserService interface {
	CreateUser(ctx context.Context, creds domain.Credentials) (*domain.User, error)
	ReadUserById(ctx context.Context, id string) (*domain.User, error)
	ReadUserByPhone(ctx context.Context, phonenumber string) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type SessionService interface {
	CreateSession(ctx context.Context, userid string) (*domain.Session, error)
	ReadSession(ctx context.Context, token string) (*domain.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

// repo layer
type UserRepository interface {
	SaveUser(ctx context.Context, user domain.User) (string, error)
	ReadUserByID(ctx context.Context, id string) (*domain.User, error)
	ReadUserByPhone(ctx context.Context, phone string) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, user domain.User) error
}

type SessionRepository interface {
	SaveSession(ctx context.Context, session domain.Session, userid string) (string, error)
	FindSessionByToken(ctx context.Context, token string) (*domain.Session, error)
	DeleteSession(ctx context.Context, token string) error
}
