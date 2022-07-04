package rds

import (
	"github.com/spf13/cobra"
)

func init() {
	instancesCmd.AddCommand(instancesListCmd)
}

// RdsCmd represents the root of the rds command
var instancesCmd = &cobra.Command{
	Use:              "instances",
	Short:            "rds instance related commands",
	PersistentPreRun: configClient,
}
