package machineIdentifier

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestReadFileOrError_NonRootNoPermission(t *testing.T) {
	// Test reading file as non-root when MustAdmin is true
	content, err := readFileOrError("/etc/machine-id", true)
	
	// Should error if not root
	// Note: In CI/test environment, we're typically not root
	if os.Getuid() != 0 {
		assert.Error(t, err)
		assert.Empty(t, content)
		assert.Contains(t, err.Error(), "must be run as root")
	}
}

func TestReadFileOrError_NonExistentFile(t *testing.T) {
	content, err := readFileOrError("/nonexistent/file/path", false)
	
	assert.Error(t, err)
	assert.Empty(t, content)
}

func TestReadFileOrError_ValidFile(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test-file")
	testContent := "test-machine-id-12345\n"
	
	err := os.WriteFile(tmpFile, []byte(testContent), 0644)
	assert.NoError(t, err)
	
	// Read the file with MustAdmin=false
	content, err := readFileOrError(tmpFile, false)
	
	assert.NoError(t, err)
	// Should strip newlines
	assert.Equal(t, "test-machine-id-12345", content)
}

func TestReadFileOrError_MultipleNewlines(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test-file-newlines")
	testContent := "line1\nline2\nline3\n"
	
	err := os.WriteFile(tmpFile, []byte(testContent), 0644)
	assert.NoError(t, err)
	
	content, err := readFileOrError(tmpFile, false)
	
	assert.NoError(t, err)
	// Should strip all newlines
	assert.Equal(t, "line1line2line3", content)
}

func TestReadMachineID_NonExistent(t *testing.T) {
	content, err := readMachineID("/nonexistent/machine-id")
	
	assert.Error(t, err)
	assert.Empty(t, content)
}

func TestReadProductUUID_NonExistent(t *testing.T) {
	content, err := readProductUUID("/nonexistent/product-uuid")
	
	assert.Error(t, err)
	assert.Empty(t, content)
}

func TestGetMachineID_InvalidPaths(t *testing.T) {
	logger := logrus.New()
	
	// Test with invalid paths
	res, err := GetMachineID(logger, "/invalid/machine-id", "/invalid/product-uuid")
	
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestGetMachineID_EmptyPaths(t *testing.T) {
	logger := logrus.New()
	
	// Test with empty paths
	res, err := GetMachineID(logger, "", "")
	
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestGetMachineID_PartialSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	machineIDPath := filepath.Join(tmpDir, "machine-id")
	
	// Create only machine-id file
	err := os.WriteFile(machineIDPath, []byte("test-machine-id"), 0644)
	assert.NoError(t, err)
	
	logger := logrus.New()
	
	// Test with valid machine-id but invalid UUID path
	res, err := GetMachineID(logger, machineIDPath, "/nonexistent/uuid")
	
	// Should still error because UUID is required
	assert.Error(t, err)
	assert.Nil(t, res)
}
