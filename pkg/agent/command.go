package agent

import (
	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/spf13/cobra"
)

// Cmd creates a new cobra command for running the agent.
// It sets up the command with a short description and a run function that calls RunAgent.
func Cmd(cfg *options.RunOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "agent",
		Short: "Run facter in agent mode",
		Run: func(cmd *cobra.Command, args []string) {
			Run(cfg)
		},
	}
}
