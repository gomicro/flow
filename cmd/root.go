package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gomicro/flow/cmd/auth"
	"github.com/gomicro/flow/fmt"
)

func init() {
	cobra.OnInitialize(initEnvs)

	RootCmd.PersistentFlags().Bool("verbose", false, "show more verbose output")
	viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))

	RootCmd.AddCommand(auth.AuthCmd)
}

func initEnvs() {
}

// RootCmd represents the base command without any params on it.
var RootCmd = &cobra.Command{
	Use:   "flow",
	Short: "A CLI for deploying services",
	Long: `Flow is a CLI tool to for a continuous deploy pipeline to use for
	deploying to cloud providers. Easily auth and roll services.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("Failed to execute: %v", err.Error())
		os.Exit(1)
	}
}
