package main

import (
	"os"
	"testing"
)

func TestCheckEnvVar(t *testing.T) {
	// Backup original environment variable value
	originalKey := os.Getenv("OPEN_AI_KEY")
	defer os.Setenv("OPEN_AI_KEY", originalKey)

	// Test with OPEN_AI_KEY set
	os.Setenv("OPEN_AI_KEY", "test_key")
	checkEnvVar()

	// Test with OPEN_AI_KEY not set
	os.Setenv("OPEN_AI_KEY", "")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("checkEnvVar() should panic when OPEN_AI_KEY is not set")
		}
	}()
	checkEnvVar()
}

func TestExtractShellCommand(t *testing.T) {
	response := "Action: Shell[ls]"
	expected := "ls"
	result, err := extractShellCommand(response)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	response = "Action: Finish[]"
	_, err = extractShellCommand(response)
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	response = ""
	_, err = extractShellCommand(response)
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}
