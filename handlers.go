package main

import (
	"context"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type AuthService interface {
	RegisterService
	LoginService
	LogoutService
	CookieService
}

type RegisterService interface {
	UserCanRegister(ctx *context.Context)

	Register(ctx *context.Context)
}

type LoginService interface {
	UserCanLogin(ctx *context.Context)
	Login(ctx *context.Context)
}

type LogoutService interface {
	UserCanLogout(ctx *context.Context)
	Logout(ctx *context.Context)
}

type CookieService interface {
	Create(ctx *context.Context)
	Read(ctx *context.Context)
	Update(ctx *context.Context)
	Delete(ctx *context.Context)
}

type AuthRepository interface {
	UserRepository
	SessionRepository
}

type UserRepository interface {
	CreateUser(ctx *context.Context)
	ReadUser(ctx *context.Context)
	UpdateUser(ctx *context.Context)
	DeleteUser(ctx *context.Context)
}

type SessionRepository interface {
	CreateSession(ctx *context.Context)
	ReadSession(ctx *context.Context)
	UpdateSession(ctx *context.Context)
	DeleteSession(ctx *context.Context)
}

type CookieRepository interface {
	CreateCookie(ctx *context.Context)
	ReadCookie(ctx *context.Context)
	UpdateCookie(ctx *context.Context)
	DeleteCookie(ctx *context.Context)
}

type User struct {
	ID          string `json:"id"`
	Phonenumber string `json:"phonenumber"`
}
