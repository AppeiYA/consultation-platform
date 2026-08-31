package redis_adapter

import (
	"context"
	"encoding/json"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/redis"
)

const MatchingQueueKey = "queue:expertmatching:jobs"

type RedisMatchingJobEnqueuer struct {
	redisClient *redis.Redis
}

func NewRedisMatchingJobEnqueuer(redisClient *redis.Redis) *RedisMatchingJobEnqueuer {
	return &RedisMatchingJobEnqueuer{
		redisClient: redisClient,
	}
}

type jobPayload struct {
	RunID  string `json:"run_id"`
	CaseID string `json:"case_id"`
}

func (e *RedisMatchingJobEnqueuer) Enqueue(ctx context.Context, job outbound.MatchingJob) error {
	payload, err := json.Marshal(jobPayload{
		RunID:  job.RunID,
		CaseID: job.CaseID,
	})
	if err != nil {
		return err
	}

	return e.redisClient.Client().RPush(ctx, MatchingQueueKey, payload).Err()
}

var _ outbound.MatchingJobEnqueuer = (*RedisMatchingJobEnqueuer)(nil)
