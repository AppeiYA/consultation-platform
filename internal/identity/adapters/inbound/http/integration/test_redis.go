package integration

import (
    "context"
    "testing"

    "github.com/AppeiYA/consultation-platform/internal/shared/config"
    shared_redis "github.com/AppeiYA/consultation-platform/internal/shared/redis"
)

func testRedis(t *testing.T) *shared_redis.Redis {
    t.Helper()

    cfg := config.SetupTestConfig()

    rdb, err := shared_redis.Connect(cfg.Redis)
    if err != nil {
        t.Fatalf("failed to connect to redis: %v", err)
    }

    cleanupRedis(t, rdb)

    t.Cleanup(func() {
        if err := rdb.Close(); err != nil {
            t.Errorf("failed to close redis: %v", err)
        }
    })

    return rdb
}

func cleanupRedis(t *testing.T, rdb *shared_redis.Redis) {
    t.Helper()

    if err := rdb.Client().FlushDB(context.Background()).Err(); err != nil {
        t.Fatalf("failed to flush redis: %v", err)
    }
}