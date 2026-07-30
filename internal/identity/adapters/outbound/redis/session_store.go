package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	shared_redis "github.com/AppeiYA/consultation-platform/internal/shared/redis"
	"github.com/redis/go-redis/v9"
)

type SessionStore struct {
	redis *shared_redis.Redis
	clock *system.SystemClock
}

func NewSessionStore(redis *shared_redis.Redis, clock *system.SystemClock) *SessionStore {
	return &SessionStore{
		redis: redis,
		clock: clock,
	}
}

func (s *SessionStore) Save(ctx context.Context, session *domain.Session) error {
	model := NewSessionModel(session)
	value, err := json.Marshal(model)
	if err != nil {
		return err
	}

	ttl := time.Until(model.ExpiresAt)
	err = s.redis.Client().Set(ctx, model.TokenHash, value, ttl).Err()

	return err
}

func (s *SessionStore) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error) {
	value, err := s.redis.Client().Get(ctx, tokenHash).Result()

	if err == redis.Nil {
		return nil, domain.ErrSessionNotFound
	}

	if err != nil {
		return nil, err
	}

	var model SessionModel

	err = json.Unmarshal([]byte(value), &model)
	if err != nil {
		return nil, err
	}
	return model.ToDomain()
}

func (s *SessionStore) Delete(ctx context.Context, tokenHash string) error {
	err := s.redis.Client().Del(ctx, tokenHash).Err()
	return err
}