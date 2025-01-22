package main

import "context"

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
