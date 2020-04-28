package ecs

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/gomicro/flow/fmt"
	"github.com/spf13/cobra"
)

func init() {
	EcsCmd.AddCommand(CreateCmd)

	UpdateCmd.Flags().StringVar(&name, "name", "", "name to create a new revision under")
}

// UpdateCmd represents the command to run a single task within ECS
var UpdateCmd = &cobra.Command{
	Use:   "create <image>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create a task definition in ECS",
	Long:  `Create a task definition for running later in ECS`,
	Run:   updateFunc,
}

func updateFunc(cmd *cobra.Command, args []string) {
	image := args[0]

	if name == "" {
		splits := strings.Split(image, "amazonaws.com/")

		imageName := splits[0]
		if len(splits) > 1 {
			imageName = splits[1]
		}

		name = imageName
	}

	descInput := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &name,
	}

	result, err := ecsSvc.DescribeTaskDefinition(descInput)
	if err != nil {
		fmt.Printf("Error describing task definition: %v", err.Error())
		os.Exit(1)
	}

	regInput := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: result.TaskDefinition.ContainerDefinitions,
	}

	regInput.ContainerDefinitions[0].Image = &image

	_, err = ecsSvc.RegisterTaskDefinition(regInput)
	if err != nil {
		fmt.Printf("Error creating task definition: %v", err.Error())
		os.Exit(1)
	}
}
