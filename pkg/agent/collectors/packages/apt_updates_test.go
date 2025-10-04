package packages

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/klamhq/facter-oss/pkg/utils"
)

func TestParseAptUpgradableOutput(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect map[string]string
	}{
		{
			name:   "single package",
			input:  `libfoo/now 1.2.3-1 amd64 [upgradable from: 1.2.2-1]`,
			expect: map[string]string{"libfoo": "1.2.3-1"},
		},
		{
			name:   "with epoch",
			input:  `mypkg/now 2:4.5.6-1 amd64 [upgradable from: 2:4.5.5-1]`,
			expect: map[string]string{"mypkg": "4.5.6-1"}, // note: regex captures version after optional epoch
		},
		{
			name:   "no matches",
			input:  `Listing...`,
			expect: map[string]string{},
		},
		{
			name: "multiple lines",
			input: `libfoo/now 1.2.3-1 amd64 [upgradable from: 1.2.2-1]
bar/now 3.4.5-2 amd64 [upgradable from: 3.4.4-2]
`,
			expect: map[string]string{"libfoo": "1.2.3-1", "bar": "3.4.5-2"},
		},
		{
			name: "ignore malformed line",
			input: `not-a-valid-line
libok/now 1.0.1 amd64 [upgradable from: 1.0.0]
`,
			expect: map[string]string{"libok": "1.0.1"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := parseAptUpgradableOutput([]byte(c.input))
			if !reflect.DeepEqual(got, c.expect) {
				t.Fatalf("unexpected result\ngot:  %#v\nwant: %#v", got, c.expect)
			}
		})
	}
}

func TestGetAptUpgradableMap_RunCmdMock(t *testing.T) {
	orig := utils.RunCmd
	defer func() { utils.RunCmd = orig }() // restore after test

	// happy path
	utils.RunCmd = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		// assert the command invoked if you want:
		if name != "apt" {
			t.Fatalf("expected apt command, got %s", name)
		}
		return []byte("libfoo/now 1.2.3-1 amd64 [upgradable from: 1.2.2-1]\n"), nil
	}
	m, err := GetAptUpgradableMap(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "1.2.3-1"; m["libfoo"] != want {
		t.Fatalf("expected libfoo version %s, got %s", want, m["libfoo"])
	}

	// error path
	utils.RunCmd = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return nil, errors.New("command failed")
	}
	_, err = GetAptUpgradableMap(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
