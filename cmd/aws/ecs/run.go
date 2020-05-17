package ecs

import (
	// "encoding/base64"
	"os"
	// "strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/gomicro/flow/fmt"
	"github.com/spf13/cobra"
)

func init() {
	EcsCmd.AddCommand(RunCmd)
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
		//Cluster:        aws.String("default"),
		//TaskDefinition: aws.String("sleep360:1"),
	}

	result, err := ecsSvc.RunTask(input)
	if err != nil {
		fmt.Printf("Error running task: %v", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%v", result)
}
