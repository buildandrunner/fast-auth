package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/mar-cial/space-auth/internal/core/domain"
	"github.com/mar-cial/space-auth/internal/core/port"
	"github.com/redis/go-redis/v9"
)

// Define all errors at the top for consistency and maintainability
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrSerialization   = errors.New("failed to serialize data")
	ErrDeserialization = errors.New("failed to deserialize data")
	ErrSaveToRedis     = errors.New("failed to save to Redis")
	ErrFetchFromRedis  = errors.New("failed to fetch from Redis")
	ErrDeleteFromRedis = errors.New("failed to delete from Redis")
)

type RedisAuthRepository struct {
	client *redis.Client
}

func NewRedisAuthRepository(client *redis.Client) port.AuthRepository {
	return &RedisAuthRepository{client: client}
}

// SaveUser saves a user in Redis.
func (r *RedisAuthRepository) SaveUser(ctx context.Context, user *domain.User) error {
	userKey := "user:" + user.Phonenumber
	userData, err := json.Marshal(user)
	if err != nil {
		return ErrSerialization
	}

	if err := r.client.Set(ctx, userKey, userData, 0).Err(); err != nil {
		return ErrSaveToRedis
	}

	return nil
}

func (r *RedisAuthRepository) FindUserByPhone(ctx context.Context, phonenumber string) (*domain.User, error) {
	userKey := "user:" + phonenumber

	// Get user data from Redis
	userData, err := r.client.Get(ctx, userKey).Result()
	if err == redis.Nil {
		// Handle "key not found" error gracefully
		return nil, ErrUserNotFound
	} else if err != nil {
		// Handle other Redis errors
		return nil, ErrFetchFromRedis
	}

	// Deserialize user data
	var user domain.User
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return nil, ErrDeserialization
	}

	return &user, nil
}

// CreateSession saves a session in Redis.
func (r *RedisAuthRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	sessionKey := "session:" + session.Token
	sessionData, err := json.Marshal(session)
	if err != nil {
		return ErrSerialization
	}

	// Save session data in Redis with an expiration time
	if err := r.client.Set(ctx, sessionKey, sessionData, time.Until(session.ExpiresAt)).Err(); err != nil {
		return ErrSaveToRedis
	}

	return nil
}

// FindSessionByToken fetches a session from Redis by token.
func (r *RedisAuthRepository) FindSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	sessionKey := "session:" + token

	// Get session data from Redis
	sessionData, err := r.client.Get(ctx, sessionKey).Result()
	if err == redis.Nil {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, ErrFetchFromRedis
	}

	// Deserialize session data
	var session domain.Session
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		return nil, ErrDeserialization
	}

	return &session, nil
}

// DeleteSession removes a session from Redis by token.
func (r *RedisAuthRepository) DeleteSession(ctx context.Context, token string) error {
	sessionKey := "session:" + token

	// Delete the session data
	if err := r.client.Del(ctx, sessionKey).Err(); err != nil {
		return ErrDeleteFromRedis
	}

	return nil
}
