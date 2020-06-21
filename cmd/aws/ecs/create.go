package ecs

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/gomicro/flow/envs"
	"github.com/gomicro/flow/fmt"
	"github.com/spf13/cobra"
)

var (
	envFile string
)

func init() {
	EcsCmd.AddCommand(CreateCmd)

	CreateCmd.Flags().StringVar(&envFile, "envFile", "", "file to parse env vars from")
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

	if envFile != "" {
		es, err := envs.ParseFile(envFile)
		if err != nil {
			fmt.Printf("Error parsing env file: %v", err.Error())
			os.Exit(1)
		}

		awsEnvs := make([]*ecs.KeyValuePair, len(es))

		for i := range es {
			awsEnvs = append(awsEnvs, &ecs.KeyValuePair{Name: &es[i].Key, Value: &es[i].Value})
		}

		input.ContainerDefinitions[0].Environment = awsEnvs
	}

	if memory != 0 {
		input.ContainerDefinitions[0].Memory = &memory
	}

	if name != "" {
		input.ContainerDefinitions[0].Name = &name
		input.Family = &name
	} else {
		splits := strings.Split(image, "amazonaws.com/")

		imageName := splits[0]
		if len(splits) > 1 {
			imageName = splits[1]
		}
		input.ContainerDefinitions[0].Name = &imageName
		input.Family = &imageName
	}

	_, err := ecsSvc.RegisterTaskDefinition(input)
	if err != nil {
		fmt.Printf("Error creating task definition: %v", err.Error())
		os.Exit(1)
	}
}
