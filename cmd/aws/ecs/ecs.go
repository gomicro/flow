package ecs

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

var (
	ecsSvc *ecs.ECS

	region string

	cpu    int64
	memory int64
	name   string
)

func init() {
	EcsCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "aws region to use")

	EcsCmd.PersistentFlags().Int64Var(&cpu, "cpu", int64(0), "cpus to assign to the task definition")
	EcsCmd.PersistentFlags().Int64Var(&memory, "memory", int64(0), "memory to assign to the task definition")
}

// EcsCmd represents the root of the ecs command
var EcsCmd = &cobra.Command{
	Use:              "ecs",
	Short:            "ECS related commands",
	PersistentPreRun: configClient,
}

func configClient(cmd *cobra.Command, args []string) {
	httpClient := &http.Client{}

	cnf := &aws.Config{
		Region:     &region,
		HTTPClient: httpClient,
	}

	sess, err := session.NewSession(cnf)
	if err != nil {
		fmt.Printf("Error creating session: %v", err.Error())
		os.Exit(1)
	}

	ecsSvc = ecs.New(sess)
}
