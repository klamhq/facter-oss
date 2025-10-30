package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// getBinaryPath returns the path to the facter binary
func getBinaryPath(t *testing.T) string {
	// First try to find the binary in the bin directory
	binPath := filepath.Join("..", "..", "bin", "facter-oss")
	if _, err := os.Stat(binPath); err == nil {
		return binPath
	}

	// If not found, build it
	t.Log("Binary not found, building...")
	buildCmd := exec.Command("make", "build")
	buildCmd.Dir = filepath.Join("..", "..")
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Check if it exists now
	if _, err := os.Stat(binPath); err != nil {
		t.Fatalf("Binary not found after build: %v", err)
	}

	return binPath
}

// TestCLIWithConfigFile tests running the CLI with a config file
func TestCLIWithConfigFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	binaryPath := getBinaryPath(t)
	tmpDir := t.TempDir()

	// Create a test config file
	configPath := filepath.Join(tmpDir, "config.yml")
	storePath := filepath.Join(tmpDir, "store.db")
	outputPath := filepath.Join(tmpDir, "output.json")

	configContent := `facter:
  enabled: true
  store:
    path: "` + storePath + `"
  logs:
    debugMode: false
  performanceProfiling:
    enabled: false
  sink:
    output:
      format: "json"
      type: "file"
      outputDirectory: "` + tmpDir + `"
      outputFilename: "output.json"
  inventory:
    platform:
      enabled: true
      hardware:
        enabled: false
      os:
        enabled: true
      kernel:
        enabled: false
      virtualization:
        enabled: false
    packages:
      enabled: false
    process:
      enabled: false
    ssh:
      enabled: false
    user:
      enabled: false
    networks:
      enabled: false
    systemdService:
      enabled: false
    applications:
      enabled: false
  compliance:
    enabled: false
  vulnerabilities:
    enabled: false
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Run the command
	cmd := exec.Command(binaryPath, "--config", configPath)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "CLI should run successfully. Output: %s", string(output))

	// Verify output file was created
	_, err = os.Stat(outputPath)
	assert.NoError(t, err, "Output file should be created")

	// Verify store was created
	_, err = os.Stat(storePath)
	assert.NoError(t, err, "Store file should be created")

	// Verify output contains expected content
	assert.Contains(t, string(output), "Collecting system facts", "Output should contain success message")
}

// TestCLIWithJSONOutput tests JSON output format
func TestCLIWithJSONOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	binaryPath := getBinaryPath(t)
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.yml")
	storePath := filepath.Join(tmpDir, "store.db")
	outputPath := filepath.Join(tmpDir, "output.json")

	configContent := `facter:
  enabled: true
  store:
    path: "` + storePath + `"
  logs:
    debugMode: false
  sink:
    output:
      format: "json"
      type: "file"
      outputDirectory: "` + tmpDir + `"
      outputFilename: "output.json"
  inventory:
    platform:
      enabled: true
      hardware:
        enabled: false
    packages:
      enabled: false
    process:
      enabled: false
    ssh:
      enabled: false
    user:
      enabled: false
    networks:
      enabled: false
    systemdService:
      enabled: false
    applications:
      enabled: false
  compliance:
    enabled: false
  vulnerabilities:
    enabled: false
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Run the command
	cmd := exec.Command(binaryPath, "--config", configPath)
	_, err = cmd.CombinedOutput()
	assert.NoError(t, err)

	// Verify JSON output is valid
	data, err := os.ReadFile(outputPath)
	assert.NoError(t, err, "Should be able to read output file")

	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	assert.NoError(t, err, "Output should be valid JSON")
	assert.NotEmpty(t, jsonData, "JSON output should not be empty")
}

// TestCLIWithProtoOutput tests protobuf output format
func TestCLIWithProtoOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	binaryPath := getBinaryPath(t)
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.yml")
	storePath := filepath.Join(tmpDir, "store.db")
	outputPath := filepath.Join(tmpDir, "output.iya")

	configContent := `facter:
  enabled: true
  store:
    path: "` + storePath + `"
  logs:
    debugMode: false
  sink:
    output:
      format: "proto"
      type: "file"
      outputDirectory: "` + tmpDir + `"
      outputFilename: "output.iya"
  inventory:
    platform:
      enabled: true
      hardware:
        enabled: false
    packages:
      enabled: false
    process:
      enabled: false
    ssh:
      enabled: false
    user:
      enabled: false
    networks:
      enabled: false
    systemdService:
      enabled: false
    applications:
      enabled: false
  compliance:
    enabled: false
  vulnerabilities:
    enabled: false
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Run the command
	cmd := exec.Command(binaryPath, "--config", configPath)
	_, err = cmd.CombinedOutput()
	assert.NoError(t, err)

	// Verify proto output file was created and has content
	data, err := os.ReadFile(outputPath)
	assert.NoError(t, err, "Should be able to read output file")
	assert.NotEmpty(t, data, "Proto output should not be empty")
}

// TestCLIMultipleRuns tests running the CLI multiple times
func TestCLIMultipleRuns(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	binaryPath := getBinaryPath(t)
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.yml")
	storePath := filepath.Join(tmpDir, "store.db")

	configContent := `facter:
  enabled: true
  store:
    path: "` + storePath + `"
  logs:
    debugMode: false
  sink:
    output:
      format: "json"
      type: "file"
      outputDirectory: "` + tmpDir + `"
      outputFilename: "output.json"
  inventory:
    platform:
      enabled: true
      hardware:
        enabled: false
    packages:
      enabled: false
    process:
      enabled: false
    ssh:
      enabled: false
    user:
      enabled: false
    networks:
      enabled: false
    systemdService:
      enabled: false
    applications:
      enabled: false
  compliance:
    enabled: false
  vulnerabilities:
    enabled: false
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// First run
	cmd1 := exec.Command(binaryPath, "--config", configPath)
	output1, err := cmd1.CombinedOutput()
	assert.NoError(t, err, "First run should succeed")
	assert.Contains(t, string(output1), "No previous inventory", "First run should compute full inventory")

	// Second run
	cmd2 := exec.Command(binaryPath, "--config", configPath)
	output2, err := cmd2.CombinedOutput()
	assert.NoError(t, err, "Second run should succeed")
	assert.Contains(t, string(output2), "Previous inventory found", "Second run should use existing store")
}

// TestCLIWithDebugMode tests debug mode
func TestCLIWithDebugMode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	binaryPath := getBinaryPath(t)
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.yml")
	storePath := filepath.Join(tmpDir, "store.db")

	configContent := `facter:
  enabled: true
  store:
    path: "` + storePath + `"
  logs:
    debugMode: true
  sink:
    output:
      format: "json"
      type: "file"
      outputDirectory: "` + tmpDir + `"
      outputFilename: "output.json"
  inventory:
    platform:
      enabled: true
      hardware:
        enabled: false
    packages:
      enabled: false
    process:
      enabled: false
    ssh:
      enabled: false
    user:
      enabled: false
    networks:
      enabled: false
    systemdService:
      enabled: false
    applications:
      enabled: false
  compliance:
    enabled: false
  vulnerabilities:
    enabled: false
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Run with debug mode
	cmd := exec.Command(binaryPath, "--config", configPath)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	// Verify debug output is present
	outputStr := string(output)
	assert.Contains(t, outputStr, "level=debug", "Debug mode should produce debug log messages")
}

// TestCLIInvalidConfig tests error handling with invalid config
func TestCLIInvalidConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	binaryPath := getBinaryPath(t)

	// Run with non-existent config file
	cmd := exec.Command(binaryPath, "--config", "/tmp/nonexistent-config.yml")
	output, err := cmd.CombinedOutput()
	assert.Error(t, err, "Should fail with non-existent config")
	assert.Contains(t, strings.ToLower(string(output)), "error", "Output should contain error message")
}

// TestCLIWithMinimalConfig tests with minimal configuration
func TestCLIWithMinimalConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	binaryPath := getBinaryPath(t)
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.yml")
	storePath := filepath.Join(tmpDir, "store.db")

	// Minimal config - only required fields
	configContent := `facter:
  enabled: true
  store:
    path: "` + storePath + `"
  sink:
    output:
      format: "json"
      type: "file"
      outputDirectory: "` + tmpDir + `"
      outputFilename: "output.json"
  inventory:
    platform:
      enabled: false
    packages:
      enabled: false
    process:
      enabled: false
    ssh:
      enabled: false
    user:
      enabled: false
    networks:
      enabled: false
    systemdService:
      enabled: false
    applications:
      enabled: false
  compliance:
    enabled: false
  vulnerabilities:
    enabled: false
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Run the command
	cmd := exec.Command(binaryPath, "--config", configPath)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "CLI should run with minimal config. Output: %s", string(output))

	// Verify store was created
	_, err = os.Stat(storePath)
	assert.NoError(t, err, "Store file should be created")
}
