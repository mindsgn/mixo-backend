package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test with non-existent env var
	result = getEnv("NON_EXISTENT_VAR", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	// Test with valid int env var
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	result := getEnvAsInt("TEST_INT", "10")
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with invalid int env var (should use default)
	os.Setenv("TEST_INVALID", "not_a_number")
	defer os.Unsetenv("TEST_INVALID")

	result = getEnvAsInt("TEST_INVALID", "10")
	if result != 10 {
		t.Errorf("Expected 10 (default), got %d", result)
	}

	// Test with non-existent env var
	result = getEnvAsInt("NON_EXISTENT_INT", "10")
	if result != 10 {
		t.Errorf("Expected 10 (default), got %d", result)
	}
}
