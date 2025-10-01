package main

import (
	"fmt"
	"os"

	"github.com/klamhq/facter-oss/pkg/agent"

	"github.com/sirupsen/logrus"

	"github.com/klamhq/facter-oss/pkg/config"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "facter",
		Short: "Facter collects system facts",
	}
	// Generate our config based on the config supplied
	// by the user in the flags
	cfgPath, err := config.ParseFlags()
	if err != nil {
		logrus.Fatal(err)
	}
	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		logrus.Fatal(err)
	}
	rootCmd.AddCommand(agent.Cmd(cfg))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
