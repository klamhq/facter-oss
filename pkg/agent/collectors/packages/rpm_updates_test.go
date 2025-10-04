package packages

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/klamhq/facter-oss/pkg/utils"
)

type mockExitError struct{ code int }

func (m *mockExitError) Error() string { return "mock exit" }
func (m *mockExitError) ExitCode() int { return m.code }

func TestGetRpmUpgradableMap(t *testing.T) {
	orig := utils.RunCmd
	defer func() { utils.RunCmd = orig }()

	tests := []struct {
		name      string
		mockFn    func(ctx context.Context, name string, args ...string) ([]byte, error)
		expect    map[string]string
		wantError bool
	}{
		{
			name: "dnf output simple",
			mockFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
				if name == "dnf" {
					return []byte("bash.x86_64 5.1.8-3.fc35 updates\nvim.noarch 9.0-1.fc35 updates\n"), nil
				}
				t.Fatalf("unexpected command: %s", name)
				return nil, nil
			},
			expect: map[string]string{
				"bash": "5.1.8-3.fc35",
				"vim":  "9.0-1.fc35",
			},
		},
		{
			name: "dnf fails, yum succeeds",
			mockFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
				if name == "dnf" {
					return nil, errors.New("dnf not found")
				}
				if name == "yum" {
					return []byte("nano.x86_64 6.0-1.fc35 updates\n"), nil
				}
				return nil, nil
			},
			expect: map[string]string{"nano": "6.0-1.fc35"},
		},
		{
			name: "no updates exit code 100 (not error)",
			mockFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
				return nil, &mockExitError{100}
			},
			expect: map[string]string{},
		},
		{
			name: "both fail with real error",
			mockFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
				return nil, errors.New("boom")
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.RunCmd = tt.mockFn
			got, err := GetRpmUpgradableMap(context.Background())
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expect) {
				t.Fatalf("expected %#v, got %#v", tt.expect, got)
			}
		})
	}
}
