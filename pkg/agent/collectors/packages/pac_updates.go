package packages

import (
	"context"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// GetPacmanUpgradableMap récupère les paquets à mettre à jour (nom -> version dispo)
func GetPacmanUpgradableMap(ctx context.Context) (map[string]string, error) {
	cmd := exec.CommandContext(ctx, "pacman", "-Qu")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=C")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// Exemple : iproute2 6.15.0-1 -> 6.16.0-1
	re := regexp.MustCompile(`^(?P<name>[^\s]+)\s+(?P<version>[^\s]+)\s+->\s+(?P<new>[^\s]+)$`)
	lines := strings.Split(string(out), "\n")
	result := make(map[string]string)

	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		if match != nil {
			name := match[1]
			newVersion := match[3]
			result[name] = newVersion
		}
	}
	return result, nil
}
