package integration

import (
	"testing"

	shared_redis "github.com/AppeiYA/consultation-platform/internal/shared/redis"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
)

func testRedis(t *testing.T) *shared_redis.Redis {
	return testhelpers.TestRedis(t)
}