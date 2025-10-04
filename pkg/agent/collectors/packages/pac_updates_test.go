package packages

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/klamhq/facter-oss/pkg/utils"
)

func TestParsePacmanOutput(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect map[string]string
	}{
		{
			name:  "single",
			input: "iproute2 6.15.0-1 -> 6.16.0-1",
			expect: map[string]string{
				"iproute2": "6.16.0-1",
			},
		},
		{
			name: "multiple",
			input: `iproute2 6.15.0-1 -> 6.16.0-1
foo-bar 1.0.0-1 -> 1.0.1-1
`,
			expect: map[string]string{
				"iproute2": "6.16.0-1",
				"foo-bar":  "1.0.1-1",
			},
		},
		{
			name:   "ignore malformed",
			input:  "not a valid line\niproute2 6.15.0-1 -> 6.16.0-1",
			expect: map[string]string{"iproute2": "6.16.0-1"},
		},
		{
			name:   "empty",
			input:  "",
			expect: map[string]string{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := parsePacmanOutput([]byte(c.input))
			if !reflect.DeepEqual(got, c.expect) {
				t.Fatalf("want %#v, got %#v", c.expect, got)
			}
		})
	}
}

func TestGetPacmanUpgradableMap_RunCmdMock(t *testing.T) {
	orig := utils.RunCmd
	defer func() { utils.RunCmd = orig }()

	// happy path
	utils.RunCmd = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		if name != "pacman" {
			t.Fatalf("expected pacman, got %s", name)
		}
		return []byte("iproute2 6.15.0-1 -> 6.16.0-1\n"), nil
	}

	m, err := GetPacmanUpgradableMap(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m["iproute2"]; got != "6.16.0-1" {
		t.Fatalf("expected iproute2 -> 6.16.0-1, got %s", got)
	}

	// error path
	utils.RunCmd = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return nil, errors.New("command failed")
	}
	_, err = GetPacmanUpgradableMap(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
