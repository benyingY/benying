package shared

import (
	"os"
	"testing"
)

func TestLoadSettingsDefaults(t *testing.T) {
	os.Unsetenv("SERVICE_HOST")
	os.Unsetenv("SERVICE_PORT")

	settings := LoadSettings("example")
	if settings.ServiceName != "example" {
		t.Fatalf("expected service name example, got %q", settings.ServiceName)
	}
	if settings.Host != "0.0.0.0" {
		t.Fatalf("expected default host 0.0.0.0, got %q", settings.Host)
	}
	if settings.Port != "8000" {
		t.Fatalf("expected default port 8000, got %q", settings.Port)
	}
}

func TestLoadSettingsOverrides(t *testing.T) {
	os.Setenv("SERVICE_HOST", "127.0.0.1")
	os.Setenv("SERVICE_PORT", "9000")
	defer os.Unsetenv("SERVICE_HOST")
	defer os.Unsetenv("SERVICE_PORT")

	settings := LoadSettings("example")
	if settings.Host != "127.0.0.1" {
		t.Fatalf("expected overridden host, got %q", settings.Host)
	}
	if settings.Port != "9000" {
		t.Fatalf("expected overridden port, got %q", settings.Port)
	}
}
