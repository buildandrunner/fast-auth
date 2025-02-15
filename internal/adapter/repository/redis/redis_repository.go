package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mar-cial/space-auth/internal/core/domain"
	"github.com/mar-cial/space-auth/internal/core/port"
	"github.com/redis/go-redis/v9"
)

var (
	baseKeyPrefix              = ""
	userKeyPrefix              = "user:"
	emailKeyPrefix             = "user:email:"
	verificationTokenKeyPrefix = "user:token:"
	sessionKeyPrefix           = "user:session:"
	accountKeyPrefix           = "user:account:"
	accountByUserIdPrefix      = "user:account:by-user-id:"
	sessionByUserIdKeyPrefix   = "user:session:by-user-id:"
	phoneByUserIdKeyPrefix     = "user:phone:by-user-id:"
)

type redisAuthRepo struct {
	client *redis.Client
}

func (r *redisAuthRepo) SaveUser(ctx context.Context, user domain.User) (string, error) {
	marshalledUser, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	// Construct Redis keys using prefixes
	userKey := userKeyPrefix + user.ID
	phoneKey := phoneByUserIdKeyPrefix + user.ID
	accountKey := accountKeyPrefix + user.ID

	// Create a pipeline for multiple commands
	pipe := r.client.TxPipeline()

	// Set the user, phone, and account details
	pipe.Set(ctx, userKey, marshalledUser, 0)
	pipe.Set(ctx, phoneKey, user.Phonenumber, 0)
	pipe.Set(ctx, accountKey, user.ID, 0)

	// Execute the pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (r *redisAuthRepo) ReadUserByID(ctx context.Context, id string) (*domain.User, error) {
	userkey := fmt.Sprintf("%s:%s", userKeyPrefix, id)

	userResponse, err := r.client.Get(ctx, userkey).Result()
	if err != nil {
		return nil, err
	}

	user := &domain.User{}
	if err := json.Unmarshal([]byte(userResponse), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *redisAuthRepo) ReadUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	userkey := fmt.Sprintf("%s:%s", phoneByUserIdKeyPrefix, phone)

	userResponse, err := r.client.Get(ctx, userkey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Return nil, nil when the user is not found
			return nil, nil
		}
		return nil, err
	}

	user := &domain.User{}
	if err := json.Unmarshal([]byte(userResponse), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *redisAuthRepo) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	// Get existing user to check for phone number changes
	existingUser, err := r.ReadUserByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	pipe := r.client.TxPipeline()

	// Update main user record
	userKey := fmt.Sprintf("user:%s", user.ID)
	userData, _ := json.Marshal(user)
	pipe.Set(ctx, userKey, userData, 0)

	// Handle phone number change
	if user.Phonenumber != existingUser.Phonenumber {
		// Delete old phone mapping
		oldPhoneKey := fmt.Sprintf("user:phone:%s", existingUser.Phonenumber)
		pipe.Del(ctx, oldPhoneKey)

		// Create new phone mapping
		newPhoneKey := fmt.Sprintf("user:phone:%s", user.Phonenumber)
		pipe.Set(ctx, newPhoneKey, userData, 0)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("update transaction failed: %w", err)
	}

	return &user, nil
}

func (r *redisAuthRepo) DeleteUser(ctx context.Context, user domain.User) error {
	pipe := r.client.TxPipeline()

	userKey := fmt.Sprintf("user:%s", user.ID)
	phoneKey := fmt.Sprintf("user:phone:%s", user.Phonenumber)

	pipe.Del(ctx, userKey)
	pipe.Del(ctx, phoneKey)

	_, err := pipe.Exec(ctx)
	return err
}

// Fix session key generation in redis/repository.go
func (r *redisAuthRepo) SaveSession(ctx context.Context, session domain.Session, userid string) (string, error) {
	sessionKey := fmt.Sprintf("user:session:%s", session.Token) // Use Token instead of ID
	sessionByIDKey := fmt.Sprintf("user:session:by-user-id:%s", userid)

	sessionBytes, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	return r.client.MSet(ctx, sessionKey, sessionBytes, sessionByIDKey, sessionBytes).Result()
}

// Add proper error type to port

// Update FindSessionByToken in redis/repository.go
func (r *redisAuthRepo) FindSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	sessionKey := fmt.Sprintf("user:session:%s", token)
	data, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, port.ErrSessionNotFound
		}
		return nil, err
	}
	session := &domain.Session{}
	if err := json.Unmarshal([]byte(data), session); err != nil {
		return nil, err
	}
	return session, nil
}

func (r *redisAuthRepo) DeleteSession(ctx context.Context, token string) error {
	sessionKey := fmt.Sprintf("user:session:%s", token)

	session, err := r.FindSessionByToken(ctx, token)
	if err != nil {
		return err
	}

	// 2. Delete both keys (adjust if session lacks UserID)
	sessionByUserKey := fmt.Sprintf("user:session:by-user-id:%s", session.UserID) // Assumes UserID exists
	if err := r.client.Del(ctx, sessionKey, sessionByUserKey).Err(); err != nil {
		return err
	}
	return nil
}

func NewRedisAuthRepository(client *redis.Client) port.AuthRepository {
	return &redisAuthRepo{client: client}
}
