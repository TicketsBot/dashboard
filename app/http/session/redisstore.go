package session

import (
	"encoding/json"
	"errors"
	"fmt"
	wrapper "github.com/TicketsBot/GoPanel/redis"
	"github.com/go-redis/redis/v8"
	"time"
)

var ErrNoSession = errors.New("no session data found")

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore() *RedisStore {
	return &RedisStore{
		client: wrapper.Client.Client,
	}
}

var keyPrefix = "panel:session:"

func (s *RedisStore) Get(userId uint64) (SessionData, error) {
	raw, err := s.client.Get(wrapper.DefaultContext(), fmt.Sprintf("%s:%d", keyPrefix, userId)).Result()
	if err != nil {
		if err == redis.Nil {
			err = ErrNoSession
		}

		return SessionData{}, err
	}

	var data SessionData
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return SessionData{}, err
	}

	return data, nil
}

func (s *RedisStore) Set(userId uint64, data SessionData) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		return err
	}

	expiration := time.Unix(data.Expiry, 0).Sub(time.Now())

	return s.client.Set(wrapper.DefaultContext(), fmt.Sprintf("%s:%d", keyPrefix, userId), encoded, expiration).Err()
}

func (s *RedisStore) Clear(userId uint64) error {
	return s.client.Del(wrapper.DefaultContext(), fmt.Sprintf("%s:%d", keyPrefix, userId)).Err()
}
