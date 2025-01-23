package main

import "context"

type UserRepository interface {
	IsPhoneBlocked(ctx context.Context, phonenumber string) (bool, error)
	CreateUser(ctx context.Context, phonenumber string) (*User, error)
	ReadUserByPhonenumber(ctx context.Context, phonenumber string) (*User, error)
	UpdateUser(ctx context.Context, phonenumber string, user User) (*User, error)
	DeleteUser(ctx context.Context, phonenumber string) error
}

type TwilioRepository interface {
	SendVerificationCode(ctx context.Context, phonenumber string) error
	CheckVerificationCode(ctx context.Context, phonenumber, code string) (bool, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context)
	ReadSession(ctx context.Context)
	DeleteSession(ctx context.Context)
}
