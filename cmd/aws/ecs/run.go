package ecs

import (
	"os"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/gomicro/flow/fmt"
	"github.com/spf13/cobra"
)

var (
	cluster string
)

func init() {
	EcsCmd.AddCommand(RunCmd)

	RunCmd.Flags().StringVar(&cluster, "cluster", "default", "the arn of the cluster to run the task on")
	RunCmd.Flags().StringVar(&name, "name", "", "name of the task definition to run")

	err := cobra.MarkFlagRequired(RunCmd.Flags(), "name")
	if err != nil {
		fmt.Printf("Error setting up: ecs: run: %v\n", err.Error())
		os.Exit(1)
	}
}

// RunCmd represents the command to run a single task within ECS
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a task in ECS",
	Long:  `Schedule a singular task to run within an ECS cluser`,
	Run:   runFunc,
}

func runFunc(cmd *cobra.Command, args []string) {
	input := &ecs.RunTaskInput{
		Cluster:        &cluster,
		TaskDefinition: &name,
	}

	_, err := ecsSvc.RunTask(input)
	if err != nil {
		fmt.Printf("Error running task: %v", err.Error())
		os.Exit(1)
	}
}
