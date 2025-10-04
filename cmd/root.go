package cmd

import (
	"github.com/klamhq/facter-oss/pkg/agent"
	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{

		Use:   "facter",
		Short: "Facter collects system facts",
		Long: `Facter is a tool that collects and reports information about
a system's hardware, operating system, and environment. It is commonly
used in configuration management systems to provide data for making
decisions about how to configure systems.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var cfg options.RunOptions
			if err := viper.Unmarshal(&cfg); err != nil {
				logrus.Fatalf("Failed to unmarshal config: %v", err)
			}
			return agent.Run(&cfg)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}

}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to facter config file")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Fatal("No config file found")
		} else {
			logrus.Info("Using config file:", viper.ConfigFileUsed())
		}

	}
}
