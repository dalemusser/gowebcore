package config

import (
	"os"
	"testing"
)

type sample struct {
	Base
	Secret string `mapstructure:"secret"`
}

func TestLoadEnv(t *testing.T) {
	os.Setenv("APP_SECRET", "xyz")
	defer os.Unsetenv("APP_SECRET")

	var cfg sample
	if err := Load(&cfg, WithEnvPrefix("APP")); err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if cfg.Secret != "xyz" {
		t.Fatalf("expected env to populate secret, got %q", cfg.Secret)
	}
}
