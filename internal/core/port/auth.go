package port

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mar-cial/space-auth/internal/core/domain"
)

type AuthHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type AuthService interface {
	Register(ctx context.Context, phonenumber string) (*domain.User, error)
	Login(ctx context.Context, phonenumber string) (*domain.User, error)
	Logout(ctx context.Context, token string) error
}

type AuthRepository interface {
	UserRepository
	SessionRepository
}

type UserRepository interface {
	SaveUser(ctx context.Context, user *domain.User) error
	FindUserByPhone(ctx context.Context, phonenumber string) (*domain.User, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	FindSessionByToken(ctx context.Context, token string) (*domain.Session, error)
	DeleteSession(ctx context.Context, token string) error
}
