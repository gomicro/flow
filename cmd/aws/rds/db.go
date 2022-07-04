package rds

import (
	"github.com/spf13/cobra"
)

func init() {
	dbCmd.AddCommand(dbInitCmd)
}

// RdsCmd represents the root of the rds command
var dbCmd = &cobra.Command{
	Use:              "db",
	Short:            "db related commands",
	PersistentPreRun: configClient,
}
