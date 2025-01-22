package main

import "context"

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
