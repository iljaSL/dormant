package cmd

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	inactivityDuration uint = 6 // months
	cfgFile            string

	rootCmd = &cobra.Command{
		Use:   "dormant",
		Short: "A tool to analyze go.mod files for inactive Go packages",
		Long: `Dormant is a CLI tool for analyzing go.mod files for inactive Go packages.
This tool should help Go developers to analyze their go.mod files and report packages
that are not being actively maintained anymore.
		`,
		Run:     root,
		Version: "0.1",
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dormant.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dormant" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dormant")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		pterm.Info.Println("Using config file:", viper.ConfigFileUsed())

		if viper.IsSet("inactivityDuration") {
			inactivityDuration = viper.GetUint("inactivityDuration")
		}
	}
}
