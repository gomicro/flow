package rds

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/gomicro/flow/client/aws/session"
	"github.com/spf13/cobra"
)

var (
	rdsSvc *rds.RDS
)

func init() {
}

// RdsCmd represents the root of the rds command
var RdsCmd = &cobra.Command{
	Use:              "rds",
	Short:            "rds related commands",
	PersistentPreRun: configClient,
}

func configClient(cmd *cobra.Command, args []string) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Printf("Error creating session: %v", err.Error())
		os.Exit(1)
	}

	rdsSvc = rds.New(sess)
}
