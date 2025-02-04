package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mar-cial/space-auth/internal/core/domain"
	"golang.org/x/crypto/argon2"
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
func generateSession(userID string, durationHours int) (*domain.Session, error) {
	if durationHours <= 0 {
		return nil, ErrInvalidSessionDuration
	}

	return &domain.Session{
		ID:        generateUniqueID(),
		Token:     generateToken(),
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(durationHours) * time.Hour),
	}, nil
}

type Argon2Params struct {
	Memory      uint32 // Memory in KiB
	Iterations  uint32 // Number of passes
	Parallelism uint8  // Number of threads
	SaltLength  uint32 // Salt length in bytes
	KeyLength   uint32 // Derived key length in bytes
}

func defaultArgon2Params() *Argon2Params {
	return &Argon2Params{
		Memory:      64 * 1024, // 64 MB
		Iterations:  1,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// Helper functions
func generateFromPassword(password string, params *Argon2Params) (string, error) {
	salt, err := generateRandomBytes(params.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Base64 encode the salt and hashed password
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: argon2id$v=19$m=65536,t=1,p=4$salt$hash
	encoded := fmt.Sprintf(
		"argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

func comparePasswordAndHash(password, encodedHash string) (bool, error) {
	// Parse encoded hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, errors.New("incompatible argon2 version")
	}

	params := &Argon2Params{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	params.SaltLength = uint32(len(salt))
	params.KeyLength = uint32(len(storedHash))

	// Derive key with same parameters
	comparisonHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Constant time comparison
	if subtle.ConstantTimeCompare(comparisonHash, storedHash) == 1 {
		return true, nil
	}
	return false, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
