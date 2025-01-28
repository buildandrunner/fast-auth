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

/*
const defaultOptions = {
  baseKeyPrefix: "",
  accountKeyPrefix: "user:account:",
  accountByUserIdPrefix: "user:account:by-user-id:",
  emailKeyPrefix: "user:email:",
  sessionKeyPrefix: "user:session:",
  sessionByUserIdKeyPrefix: "user:session:by-user-id:",
  userKeyPrefix: "user:",
  verificationTokenKeyPrefix: "user:token:",
}
*/

type redisAuthRepo struct {
	client *redis.Client
}

func (r *redisAuthRepo) SaveUser(ctx context.Context, user domain.User) (string, error) {
	userKey := fmt.Sprintf("user:%s", user.ID)
	userPhoneKey := fmt.Sprintf("user:phone:%s", user.Phonenumber)

	marshalledUser, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	return r.client.MSet(ctx, userKey, marshalledUser, userPhoneKey, user.Phonenumber).Result()
}

func (r *redisAuthRepo) ReadUserByID(ctx context.Context, id string) (*domain.User, error) {
	userkey := fmt.Sprintf("user:%s", id)

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
	userkey := fmt.Sprintf("user:phone:%s", phone)

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

func (r *redisAuthRepo) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	userkey := fmt.Sprintf("user:%s", user.ID)

	userResponse, err := r.client.Get(ctx, userkey).Result()
	if err != nil {
		return nil, err
	}

	u := &domain.User{}
	if err := json.Unmarshal([]byte(userResponse), u); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *redisAuthRepo) DeleteUser(ctx context.Context, user domain.User) error {
	userkey := fmt.Sprintf("user:%s", user.ID)

	deleteUserResponse, err := r.client.Del(ctx, userkey).Result()
	if err != nil {
		return err
	}

	if deleteUserResponse != int64(1) {
		return errors.New("err deleting user")
	}

	return nil
}

func (r *redisAuthRepo) SaveSession(ctx context.Context, session domain.Session, userid string) (string, error) {
	sessionkey := fmt.Sprintf("user:session:%s", session.ID)
	sessionByIDkey := fmt.Sprintf("user:session:by-user-id:%s", userid)

	sessionbytes, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	return r.client.MSet(ctx, sessionkey, sessionbytes, sessionByIDkey, sessionbytes).Result()
}

func (r *redisAuthRepo) FindSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	sessionKey := fmt.Sprintf("user:session:%s", token)
	data, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("session not found")
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
