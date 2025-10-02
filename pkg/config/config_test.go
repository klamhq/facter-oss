package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")
	yamlContent := `
field1: "value1"
field2: 42
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	config, err := NewConfig(configPath)
	assert.NoError(t, err)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assert.NotNil(t, config)
	// Optionally check fields if RunOptions has exported fields
}

func TestNewConfig_FileNotFound(t *testing.T) {
	_, err := NewConfig("nonexistent.yml")
	assert.Error(t, err)
}

func TestNewConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "bad.yml")
	badYAML := "a: [unclosed"
	if err := os.WriteFile(configPath, []byte(badYAML), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	_, err := NewConfig(configPath)
	assert.Error(t, err)
}

func TestValidateConfigPath_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")
	if err := os.WriteFile(configPath, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	err := ValidateConfigPath(configPath)
	assert.NoError(t, err)
}

func TestValidateConfigPath_NonExistent(t *testing.T) {
	err := ValidateConfigPath("nonexistent.yml")
	assert.Error(t, err)
}

func TestValidateConfigPath_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	err := ValidateConfigPath(tmpDir)
	assert.Error(t, err)
}

func TestParseFlags(t *testing.T) {
	// Create a temp config file to point to
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")
	if err := os.WriteFile(configPath, []byte("field: value"), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	// Simulate command line argument
	os.Args = []string{"cmd", "--config", configPath}

	returnedPath, err := ParseFlags()
	assert.NoError(t, err)
	assert.Equal(t, configPath, returnedPath)
}

func TestParseFlags_InvalidPath(t *testing.T) {
	// Simulate command line argument with invalid path
	os.Args = []string{"cmd", "--config", "nonexistent.yml"}

	_, err := ParseFlags()
	assert.Error(t, err)
}

func TestParseFlags_DirectoryPath(t *testing.T) {
	// Create a temp directory to simulate directory path
	tmpDir := t.TempDir()

	// Simulate command line argument with directory path
	os.Args = []string{"cmd", "--config", tmpDir}

	_, err := ParseFlags()
	assert.Error(t, err)
}
