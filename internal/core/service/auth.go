package service

import (
	"context"

	"github.com/mar-cial/space-auth/internal/core/domain"
	"github.com/mar-cial/space-auth/internal/core/port"
)

type authService struct {
	authRepo port.AuthRepository
}

func (a *authService) ValidateUser(ctx context.Context, creds domain.Credentials) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func (a *authService) CreateUser(ctx context.Context, creds domain.Credentials) (*domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (a *authService) ReadUserById(ctx context.Context, id string) (*domain.User, error) {
	panic("not implemented") // TODO: Implement
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
