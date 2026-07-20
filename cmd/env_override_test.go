package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// TestEnvVarOverridesOutputType verifies that env vars override nested keys.
// Convention: facter.sink.output.type → FACTER_SINK_OUTPUT_TYPE
func TestEnvVarOverridesOutputType(t *testing.T) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("facter.sink.output.type", "file")

	t.Setenv("FACTER_SINK_OUTPUT_TYPE", "remote")

	got := v.GetString("facter.sink.output.type")
	if got != "remote" {
		t.Errorf("expected %q, got %q", "remote", got)
	}
}

func TestEnvVarOverridesGrpcHost(t *testing.T) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("facter.sink.output.facterServer.serverHost", "localhost")

	t.Setenv("FACTER_SINK_OUTPUT_FACTERSERVER_SERVERHOST", "grpc.production.internal")

	got := v.GetString("facter.sink.output.facterServer.serverHost")
	if got != "grpc.production.internal" {
		t.Errorf("expected %q, got %q", "grpc.production.internal", got)
	}
}

func TestEnvVarOverridesInsecureSkipTLS(t *testing.T) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("facter.sink.output.facterServer.insecureSkipTlsVerify", false)

	t.Setenv("FACTER_SINK_OUTPUT_FACTERSERVER_INSECURESKIPTLSVERIFY", "true")

	got := v.GetBool("facter.sink.output.facterServer.insecureSkipTlsVerify")
	if !got {
		t.Errorf("expected insecureSkipTlsVerify to be true via env var, got false")
	}
}
