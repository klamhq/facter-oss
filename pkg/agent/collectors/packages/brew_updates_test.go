package packages

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/klamhq/facter-oss/pkg/utils"
)

func TestGetBrewUpgradableMap(t *testing.T) {
	orig := utils.RunCmd
	defer func() { utils.RunCmd = orig }()

	tests := []struct {
		name     string
		output   string
		runErr   error
		expected map[string]string
		wantErr  bool
	}{
		{
			name: "single package",
			output: `zeromq (4.3.4) < 4.3.5_2
`,
			expected: map[string]string{"zeromq": "4.3.5_2"},
		},
		{
			name: "multiple packages",
			output: `zeromq (4.3.4) < 4.3.5_2
go (1.22.0) < 1.22.2
`,
			expected: map[string]string{"zeromq": "4.3.5_2", "go": "1.22.2"},
		},
		{
			name:    "error case",
			runErr:  errors.New("command failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.RunCmd = func(ctx context.Context, name string, args ...string) ([]byte, error) {
				if tt.runErr != nil {
					return nil, tt.runErr
				}
				return []byte(tt.output), nil
			}

			got, err := GetBrewUpgradableMap(context.Background())
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("expected %#v, got %#v", tt.expected, got)
			}
		})
	}
}
