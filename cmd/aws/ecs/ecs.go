package ecs

import (
	"os"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/client/aws/session"
	"github.com/gomicro/flow/fmt"
)

var (
	ecsSvc *ecs.ECS

	cpu    int64
	memory int64
	name   string
)

func init() {
}

// EcsCmd represents the root of the ecs command
var EcsCmd = &cobra.Command{
	Use:              "ecs",
	Short:            "ECS related commands",
	PersistentPreRun: configClient,
}

func configClient(cmd *cobra.Command, args []string) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Printf("Error creating session: %v", err.Error())
		os.Exit(1)
	}

	ecsSvc = ecs.New(sess)
}
