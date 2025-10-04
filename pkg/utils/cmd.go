package utils

import (
	"context"
	"os"
	"os/exec"
)

var RunCmd = func(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=C")
	return cmd.CombinedOutput()
}
