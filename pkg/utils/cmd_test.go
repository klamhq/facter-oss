package utils

import (
	"context"
	"testing"
)

// Mock de RunCmd pour test
func TestRunCmd_Mock(t *testing.T) {
	orig := RunCmd
	defer func() { RunCmd = orig }()

	RunCmd = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		if name != "mycmd" {
			t.Fatalf("expected mycmd, got %s", name)
		}
		if len(args) != 2 || args[0] != "arg1" || args[1] != "arg2" {
			t.Fatalf("unexpected args: %v", args)
		}
		return []byte("mock output"), nil
	}

	out, err := RunCmd(context.Background(), "mycmd", "arg1", "arg2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "mock output" {
		t.Fatalf("expected mock output, got %q", string(out))
	}
}
