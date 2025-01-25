package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mar-cial/space-auth/internal/core/domain"
)

var (
	ErrBadToken               = errors.New("failed to generate secure random bytes")
	ErrInvalidSessionDuration = errors.New("invalid session duration")
)

func generateUniqueID() string {
	return uuid.NewString()
}

func generateToken() string {
	tokenBytes := make([]byte, 32)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return token

}

// CreateSession generates a new session for a user.
func createSession(userID string, durationHours int) (*domain.Session, error) {
	if durationHours <= 0 {
		return nil, ErrInvalidSessionDuration
	}

	return &domain.Session{
		ID:        generateUniqueID(), // Unique session ID
		Token:     generateToken(),    // Secure token (32 bytes)
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(durationHours) * time.Hour),
	}, nil
}
