package ecs

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/gomicro/flow/fmt"
	"github.com/spf13/cobra"
)

var (
	cpu    int64
	memory int64
	name   string
)

func init() {
	EcsCmd.AddCommand(CreateCmd)

	CreateCmd.Flags().Int64Var(&cpu, "cpu", int64(0), "cpus to assign to the task definition")
	CreateCmd.Flags().Int64Var(&memory, "memory", int64(0), "memory to assign to the task definition")
	CreateCmd.Flags().StringVar(&name, "name", "", "name to give the task definition")
}

// CreateCmd represents the command to run a single task within ECS
var CreateCmd = &cobra.Command{
	Use:   "create <image> [command]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create a task definition in ECS",
	Long:  `Create a task definition for running later in ECS`,
	Run:   createFunc,
}

func createFunc(cmd *cobra.Command, args []string) {
	image := args[0]

	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Image: &image,
			},
		},
	}

	command := []*string{}
	if len(args) > 1 {
		args = args[1:]
		for i := range args {
			command = append(command, &args[i])
		}
	}

	if len(command) > 0 {
		input.ContainerDefinitions[0].Command = command
	}

	if cpu != 0 {
		input.ContainerDefinitions[0].Cpu = &cpu
	}

	if memory != 0 {
		input.ContainerDefinitions[0].Memory = &memory
	}

	if name != "" {
		input.ContainerDefinitions[0].Name = &name
	} else {
		splits := strings.Split(image, "amazonaws.com/")

		imageName := splits[0]
		if len(splits) > 1 {
			imageName = splits[1]
		}
		input.ContainerDefinitions[0].Name = &imageName
	}

	splits := strings.Split(image, "/")
	family := splits[len(splits)-1]
	input.Family = &family
}
