package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var (
	// Version is the current version of flow, made available for use through
	// out the application.
	Version string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version",
	Long:  `Display the version of the CLI.`,
	Run:   versionFunc,
}

func versionFunc(cmd *cobra.Command, args []string) {
	if Version == "" {
		fmt.Printf("Flow version dev-local")
	} else {
		fmt.Printf("Flow version %v", Version)
	}
}
