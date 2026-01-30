package config

import "testing"

func TestLoadConfigDev(t *testing.T) {
    t.Setenv("APP_ENV", "dev")
    cfg := Load()
    if cfg.Env != "dev" {
        t.Fatalf("expected env dev, got %s", cfg.Env)
    }
    if cfg.Greeting != "hello-dev" {
        t.Fatalf("expected greeting hello-dev, got %s", cfg.Greeting)
    }
}

func TestLoadConfigStaging(t *testing.T) {
    t.Setenv("APP_ENV", "staging")
    cfg := Load()
    if cfg.Env != "staging" {
        t.Fatalf("expected env staging, got %s", cfg.Env)
    }
    if cfg.Greeting != "hello-staging" {
        t.Fatalf("expected greeting hello-staging, got %s", cfg.Greeting)
    }
}

func TestEnvOverride(t *testing.T) {
    t.Setenv("APP_ENV", "dev")
    t.Setenv("GREETING", "override")
    cfg := Load()
    if cfg.Greeting != "override" {
        t.Fatalf("expected greeting override, got %s", cfg.Greeting)
    }
}
