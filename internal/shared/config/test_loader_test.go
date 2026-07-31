package config

import (
	"testing"
)

func TestSetupTestConfig_Defaults(t *testing.T) {
	cfg := SetupTestConfig()

	if cfg.App.Env != "test" {
		t.Errorf("expected App.Env to be 'test', got %s", cfg.App.Env)
	}

	if cfg.Database.Name != "consultation_platform_test" {
		t.Errorf("expected Database.Name to be 'consultation_platform_test', got %s", cfg.Database.Name)
	}

	if cfg.Redis.DB != 15 {
		t.Errorf("expected Redis.DB to be 15, got %d", cfg.Redis.DB)
	}

	if cfg.Session.CookieName != "test_session_id" {
		t.Errorf("expected Session.CookieName to be 'test_session_id', got %s", cfg.Session.CookieName)
	}
}

func TestSetupTestConfig_FunctionalOptions(t *testing.T) {
	cfg := SetupTestConfig(
		WithTestDatabaseName("custom_db_test"),
		WithTestRedisDB(14),
	)

	if cfg.Database.Name != "custom_db_test" {
		t.Errorf("expected Database.Name 'custom_db_test', got %s", cfg.Database.Name)
	}

	if cfg.Redis.DB != 14 {
		t.Errorf("expected Redis.DB 14, got %d", cfg.Redis.DB)
	}
}

func TestValidateTestSafety_RejectsNonTestEnv(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic when running test config with non-test env")
		}
	}()

	_ = SetupTestConfig(func(c *Config) {
		c.App.Env = "development"
	})
}

func TestValidateTestSafety_RejectsProductionDatabase(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic when running test config with non-test database name")
		}
	}()

	_ = SetupTestConfig(func(c *Config) {
		c.Database.Name = "consultation_platform"
	})
}

func TestValidateTestSafety_RejectsRedisDBZero(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic when running test config with Redis DB 0")
		}
	}()

	_ = SetupTestConfig(func(c *Config) {
		c.Redis.DB = 0
	})
}

func TestValidateTestSafety_RejectsNonTestCookieName(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic when running test config with non-test cookie name")
		}
	}()

	_ = SetupTestConfig(func(c *Config) {
		c.Session.CookieName = "session_id"
	})
}
