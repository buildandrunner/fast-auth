package main

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrUnauthorized            = errors.New("user unauthorized")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
	ErrInvalidPhoneNumber      = errors.New("invalid phone number format")
	ErrPhoneNumberBlocked      = errors.New("phone number is blocked")
	ErrSendLimitExceeded       = errors.New("too many verification attempts")
)

// 1. design

type AuthService interface {
	UserService
	CookieService
}

type UserService interface {
	RegisterService
	LoginService
	LogoutService
}

type RegisterService interface {
	BeginPhoneVerification(ctx context.Context, phonenumber string) error
	CompletePhoneVerification(ctx context.Context, phonenumber, code string) (*User, error)
}

type LoginService interface {
	Login(ctx context.Context, phonenumber string) (*User, error)
}

type LogoutService interface {
	Logout(ctx context.Context, user User) error
}

type CookieService interface {
	CreateSessionCookie(ctx context.Context, user User) error
	ReadSesssionCookie(ctx context.Context, user User) (*http.Cookie, error)
	DeleteSessionCookie(ctx context.Context, user User) error
}

type authService struct {
	userRepo    UserRepository
	twilioRepo  TwilioRepository
	sessionRepo SessionRepository
}

func NewAuthService(
	ur UserRepository,
	tr TwilioRepository,
	sr SessionRepository,
) AuthService {
	return &authService{
		userRepo:    ur,
		twilioRepo:  tr,
		sessionRepo: sr,
	}

}

// isValidPhoneNumber is a placeholder for your phone validation logic.
func isValidPhoneNumber(phone string) bool {
	// For example: use regex, or a third-party lib like "github.com/nyaruka/phonenumbers"
	return len(phone) > 0
}

// IsPhoneBlocked might live on the UserRepository or a separate BlocklistRepository.
func (r *twilioRepo) IsPhoneBlocked(ctx context.Context, phonenumber string) (bool, error) {
	// e.g., check a "blocked_phones" table/Redis set. Return true if found.
	return false, nil // not blocked by default
}

// checkSendLimit is a placeholder for rate-limiting logic.
func (a *authService) checkSendLimit(ctx context.Context, phonenumber string) error {
	// e.g., check Redis or memory to see how many times we've sent
	// in the last X minutes. If too many, return ErrSendLimitExceeded.
	return nil
}

// logSendAttempt might be for analytics or logging.
func (a *authService) logSendAttempt(ctx context.Context, phonenumber string) error {
	// e.g., insert a "verification attempt" log into a DB, or simply log to console.
	return nil
}

func (a *authService) BeginPhoneVerification(ctx context.Context, phonenumber string) error {
	// 1) Validate phone number format (placeholder function).
	if !isValidPhoneNumber(phonenumber) {
		return ErrInvalidPhoneNumber
	}

	// 2) Optionally check if the phone number is blocked or flagged.
	//    This might be in your user repository or a separate "blocklist" repo.
	blocked, err := a.userRepo.IsPhoneBlocked(ctx, phonenumber)
	if err != nil {
		return err
	}
	if blocked {
		return ErrPhoneNumberBlocked
	}

	// 3) (Optional) Check if user or phone number is hitting rate limits (e.g., too many requests).
	//    This could be a call to Redis or an in-memory counter.
	if err := a.checkSendLimit(ctx, phonenumber); err != nil {
		return err
	}

	// 4) Send the verification code via Twilio.
	if err := a.twilioRepo.SendVerificationCode(ctx, phonenumber); err != nil {
		return err
	}

	// 5) (Optional) Record that we sent a code, for logging or analytics.
	if err := a.logSendAttempt(ctx, phonenumber); err != nil {
		// Typically we wouldn't fail the entire request if logging fails,
		// but handle it however makes sense for your use case.
	}

	// 6) Return success if no errors occurred.
	return nil
}

func (a *authService) CompletePhoneVerification(ctx context.Context, phonenumber, code string) (*User, error) {
	// 1) Check the code with Twilio
	isValid, err := a.twilioRepo.CheckVerificationCode(ctx, phonenumber, code)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, ErrInvalidVerificationCode
	}

	// 2) Check if user exists; if not, create new user
	user, err := a.userRepo.ReadUserByPhonenumber(ctx, phonenumber)
	if err != nil {
		return nil, err
	}
	if user == nil {
		newUser, err := a.userRepo.CreateUser(ctx, phonenumber)
		if err != nil {
			return nil, err
		}
		user = newUser
	}

	// 3) Return user (registered & verified!)
	return user, nil
}

func (a *authService) CheckVerificationCode(ctx context.Context, phonenumber string, code string) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func (a *authService) Login(ctx context.Context, phonenumber string) (*User, error) {
	return nil, errors.New("Login not implemented")
}

func (a *authService) Logout(ctx context.Context, user User) error {
	return errors.New("Logout not implemented")
}

func (a *authService) CreateSessionCookie(ctx context.Context, user User) error {
	return errors.New("CreateSessionCookie not implemented")
}

func (a *authService) ReadSesssionCookie(ctx context.Context, user User) (*http.Cookie, error) {
	return nil, errors.New("ReadSesssionCookie not implemented")
}

func (a *authService) DeleteSessionCookie(ctx context.Context, user User) error {
	return errors.New("DeleteSessionCookie not implemented")
}
