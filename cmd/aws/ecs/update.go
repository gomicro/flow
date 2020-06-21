package ecs

import (
	"os"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/gomicro/flow/envs"
	"github.com/gomicro/flow/fmt"
	"github.com/spf13/cobra"
)

var (
	env   []string
	image string
)

func init() {
	EcsCmd.AddCommand(UpdateCmd)

	UpdateCmd.Flags().StringSliceVarP(&env, "env", "e", []string{}, "env var key value pair to update")
	UpdateCmd.Flags().Int64Var(&cpu, "cpu", int64(0), "cpus to assign to the task definition")
	UpdateCmd.Flags().Int64Var(&memory, "memory", int64(0), "memory to assign to the task definition")
	UpdateCmd.Flags().StringVar(&image, "image", "", "image to update")
	UpdateCmd.Flags().StringVar(&name, "name", "", "name of the task definition to update and create a new revision under")
	cobra.MarkFlagRequired(UpdateCmd.Flags(), "name")
}

// UpdateCmd represents the command to run a single task within ECS
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a task definition in ECS",
	Long:  `Update a task definition in ECS with new values to run`,
	Run:   updateFunc,
}

func updateFunc(cmd *cobra.Command, args []string) {
	descInput := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &name,
	}

	result, err := ecsSvc.DescribeTaskDefinition(descInput)
	if err != nil {
		fmt.Printf("Error describing task definition: %v", err.Error())
		os.Exit(1)
	}

	// Until a usecase exists, assume only one container def
	updatedDefs := result.TaskDefinition.ContainerDefinitions

	newCPU := *updatedDefs[0].Cpu
	if cpu != 0 {
		newCPU = cpu
	}

	newMemory := *updatedDefs[0].Memory
	if memory != 0 {
		newMemory = memory
	}

	newImage := *updatedDefs[0].Image
	if image != "" {
		newImage = image
	}

	parsed := envs.ParseSlice(env)
	newEnvs := make(map[string]string)
	for i := range parsed {
		newEnvs[parsed[i].Key] = parsed[i].Value
	}

	mergedEnvs := updatedDefs[0].Environment
	for i := range mergedEnvs {
		n := *mergedEnvs[i].Name
		v, ok := newEnvs[n]
		if ok {
			mergedEnvs[i].Value = &v
		}
	}

	updatedDefs[0].Cpu = &newCPU
	updatedDefs[0].Memory = &newMemory
	updatedDefs[0].Image = &newImage
	updatedDefs[0].Environment = mergedEnvs

	regInput := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    updatedDefs,
		Cpu:                     result.TaskDefinition.Cpu,
		ExecutionRoleArn:        result.TaskDefinition.ExecutionRoleArn,
		Family:                  result.TaskDefinition.Family,
		InferenceAccelerators:   result.TaskDefinition.InferenceAccelerators,
		IpcMode:                 result.TaskDefinition.IpcMode,
		Memory:                  result.TaskDefinition.Memory,
		NetworkMode:             result.TaskDefinition.NetworkMode,
		PidMode:                 result.TaskDefinition.PidMode,
		PlacementConstraints:    result.TaskDefinition.PlacementConstraints,
		ProxyConfiguration:      result.TaskDefinition.ProxyConfiguration,
		RequiresCompatibilities: result.TaskDefinition.RequiresCompatibilities,
		TaskRoleArn:             result.TaskDefinition.TaskRoleArn,
		Volumes:                 result.TaskDefinition.Volumes,
	}

	_, err = ecsSvc.RegisterTaskDefinition(regInput)
	if err != nil {
		fmt.Printf("Error creating task definition: %v", err.Error())
		os.Exit(1)
	}
}
