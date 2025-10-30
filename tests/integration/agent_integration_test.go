package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/klamhq/facter-oss/pkg/agent"
	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/stretchr/testify/assert"
)

// TestAgentRunIntegration tests the full agent run workflow
func TestAgentRunIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")
	outputFile := filepath.Join(tmpDir, "output.json")

	opts := &options.RunOptions{}
	opts.Facter.Enabled = true
	opts.Facter.Store.Path = storePath
	opts.Facter.Logs.DebugMode = false
	opts.Facter.PerformanceProfiling.Enabled = false

	// Configure sink
	opts.Facter.Sink.Output.Format = "json"
	opts.Facter.Sink.Output.Type = "file"
	opts.Facter.Sink.Output.OutputDirectory = tmpDir
	opts.Facter.Sink.Output.OutputFilename = "output.json"

	// Enable basic collectors
	opts.Facter.Inventory.Platform.Enabled = true
	opts.Facter.Inventory.Platform.Hardware.Enabled = true
	opts.Facter.Inventory.Platform.Kernel.Enabled = true
	opts.Facter.Inventory.Platform.Os.Enabled = true

	// Disable resource-intensive collectors for faster testing
	opts.Facter.Inventory.Packages.Enabled = false
	opts.Facter.Inventory.Process.Enabled = false
	opts.Facter.Inventory.Networks.Enabled = false
	opts.Facter.Inventory.SSH.Enabled = false
	opts.Facter.Inventory.User.Enabled = false
	opts.Facter.Inventory.SystemdService.Enabled = false
	opts.Facter.Inventory.Applications.Enabled = false

	// Disable external services
	opts.Facter.Compliance.Enabled = false
	opts.Facter.Vulnerabilities.Enabled = false

	// Run the agent
	err := agent.Run(opts)
	assert.NoError(t, err, "Agent run should complete successfully")

	// Verify store was created
	_, err = os.Stat(storePath)
	assert.NoError(t, err, "Store file should be created")

	// Verify output file was created
	_, err = os.Stat(outputFile)
	assert.NoError(t, err, "Output file should be created")

	// Verify output file has content
	data, err := os.ReadFile(outputFile)
	assert.NoError(t, err, "Should be able to read output file")
	assert.NotEmpty(t, data, "Output file should have content")

	// Verify JSON is valid (if format is JSON)
	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	assert.NoError(t, err, "Output should be valid JSON")
}

// TestAgentRunWithStoreUpdate tests that the agent correctly handles incremental updates
func TestAgentRunWithStoreUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")

	opts := &options.RunOptions{}
	opts.Facter.Enabled = true
	opts.Facter.Store.Path = storePath
	opts.Facter.Logs.DebugMode = false

	// Configure sink
	opts.Facter.Sink.Output.Format = "json"
	opts.Facter.Sink.Output.Type = "file"
	opts.Facter.Sink.Output.OutputDirectory = tmpDir
	opts.Facter.Sink.Output.OutputFilename = "output.json"

	// Enable minimal collectors
	opts.Facter.Inventory.Platform.Enabled = true
	opts.Facter.Inventory.Platform.Hardware.Enabled = false
	opts.Facter.Inventory.Packages.Enabled = false
	opts.Facter.Inventory.Process.Enabled = false
	opts.Facter.Inventory.Networks.Enabled = false
	opts.Facter.Inventory.SSH.Enabled = false
	opts.Facter.Inventory.User.Enabled = false
	opts.Facter.Inventory.SystemdService.Enabled = false
	opts.Facter.Inventory.Applications.Enabled = false
	opts.Facter.Compliance.Enabled = false
	opts.Facter.Vulnerabilities.Enabled = false

	// First run - should create full inventory
	err := agent.Run(opts)
	assert.NoError(t, err, "First agent run should complete successfully")

	// Verify store exists
	info1, err := os.Stat(storePath)
	assert.NoError(t, err, "Store file should exist after first run")

	// Second run - should use existing store and compute delta
	err = agent.Run(opts)
	assert.NoError(t, err, "Second agent run should complete successfully")

	// Verify store still exists
	info2, err := os.Stat(storePath)
	assert.NoError(t, err, "Store file should still exist after second run")

	// Store should have been accessed (modified time may or may not change)
	assert.NotNil(t, info1)
	assert.NotNil(t, info2)
}

// TestAgentRunWithDifferentOutputFormats tests different output formats
func TestAgentRunWithDifferentOutputFormats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	formats := []string{"json", "proto"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			tmpDir := t.TempDir()
			storePath := filepath.Join(tmpDir, "test-store.db")
			outputExt := format
			if format == "proto" {
				outputExt = "iya"
			}
			outputFile := filepath.Join(tmpDir, "output."+outputExt)

			opts := &options.RunOptions{}
			opts.Facter.Enabled = true
			opts.Facter.Store.Path = storePath
			opts.Facter.Logs.DebugMode = false

			// Configure sink with specific format
			opts.Facter.Sink.Output.Format = format
			opts.Facter.Sink.Output.Type = "file"
			opts.Facter.Sink.Output.OutputDirectory = tmpDir
			opts.Facter.Sink.Output.OutputFilename = "output." + outputExt

			// Enable minimal collectors
			opts.Facter.Inventory.Platform.Enabled = true
			opts.Facter.Inventory.Platform.Hardware.Enabled = false
			opts.Facter.Inventory.Packages.Enabled = false
			opts.Facter.Inventory.Process.Enabled = false
			opts.Facter.Inventory.Networks.Enabled = false
			opts.Facter.Inventory.SSH.Enabled = false
			opts.Facter.Inventory.User.Enabled = false
			opts.Facter.Inventory.SystemdService.Enabled = false
			opts.Facter.Inventory.Applications.Enabled = false
			opts.Facter.Compliance.Enabled = false
			opts.Facter.Vulnerabilities.Enabled = false

			// Run the agent
			err := agent.Run(opts)
			assert.NoError(t, err, "Agent run should complete successfully for format: "+format)

			// Verify output file was created
			_, err = os.Stat(outputFile)
			assert.NoError(t, err, "Output file should be created for format: "+format)

			// Verify output file has content
			data, err := os.ReadFile(outputFile)
			assert.NoError(t, err, "Should be able to read output file for format: "+format)
			assert.NotEmpty(t, data, "Output file should have content for format: "+format)
		})
	}
}

// TestAgentRunWithPerformanceProfiling tests performance profiling mode
func TestAgentRunWithPerformanceProfiling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")

	opts := &options.RunOptions{}
	opts.Facter.Enabled = true
	opts.Facter.Store.Path = storePath
	opts.Facter.Logs.DebugMode = true
	opts.Facter.PerformanceProfiling.Enabled = true

	// Configure sink
	opts.Facter.Sink.Output.Format = "json"
	opts.Facter.Sink.Output.Type = "file"
	opts.Facter.Sink.Output.OutputDirectory = tmpDir
	opts.Facter.Sink.Output.OutputFilename = "output.json"

	// Enable minimal collectors
	opts.Facter.Inventory.Platform.Enabled = true
	opts.Facter.Inventory.Platform.Hardware.Enabled = false
	opts.Facter.Inventory.Packages.Enabled = false
	opts.Facter.Inventory.Process.Enabled = false
	opts.Facter.Inventory.Networks.Enabled = false
	opts.Facter.Inventory.SSH.Enabled = false
	opts.Facter.Inventory.User.Enabled = false
	opts.Facter.Inventory.SystemdService.Enabled = false
	opts.Facter.Inventory.Applications.Enabled = false
	opts.Facter.Compliance.Enabled = false
	opts.Facter.Vulnerabilities.Enabled = false

	// Run the agent with profiling
	err := agent.Run(opts)
	assert.NoError(t, err, "Agent run with profiling should complete successfully")

	// Note: CPU and memory profile files are created in the current directory
	// and require cleanup, but we're running in a temp environment
}

// TestAgentRunStoreCorruption tests handling of corrupted store
func TestAgentRunStoreCorruption(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")

	opts := &options.RunOptions{}
	opts.Facter.Enabled = true
	opts.Facter.Store.Path = storePath
	opts.Facter.Logs.DebugMode = false

	// Configure sink
	opts.Facter.Sink.Output.Format = "json"
	opts.Facter.Sink.Output.Type = "file"
	opts.Facter.Sink.Output.OutputDirectory = tmpDir
	opts.Facter.Sink.Output.OutputFilename = "output.json"

	// Enable minimal collectors
	opts.Facter.Inventory.Platform.Enabled = true
	opts.Facter.Inventory.Platform.Hardware.Enabled = false
	opts.Facter.Inventory.Packages.Enabled = false
	opts.Facter.Inventory.Process.Enabled = false
	opts.Facter.Inventory.Networks.Enabled = false
	opts.Facter.Inventory.SSH.Enabled = false
	opts.Facter.Inventory.User.Enabled = false
	opts.Facter.Inventory.SystemdService.Enabled = false
	opts.Facter.Inventory.Applications.Enabled = false
	opts.Facter.Compliance.Enabled = false
	opts.Facter.Vulnerabilities.Enabled = false

	// Create a corrupted store file
	err := os.WriteFile(storePath, []byte("corrupted data"), 0644)
	assert.NoError(t, err)

	// Run should handle corrupted store gracefully
	// The important part is it shouldn't panic
	assert.NotPanics(t, func() {
		_ = agent.Run(opts)
	})
}
