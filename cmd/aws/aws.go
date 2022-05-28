package aws

import (
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/cmd/aws/asm"
	"github.com/gomicro/flow/cmd/aws/ecr"
	"github.com/gomicro/flow/cmd/aws/ecs"
	"github.com/gomicro/flow/cmd/aws/s3"
)

func init() {
	AwsCmd.AddCommand(asm.AsmCmd)
	AwsCmd.AddCommand(ecr.EcrCmd)
	AwsCmd.AddCommand(ecs.EcsCmd)
	AwsCmd.AddCommand(s3.S3Cmd)
}

// AwsCmd represents the root of the aws command
var AwsCmd = &cobra.Command{
	Use:   "aws",
	Short: "AWS related commands",
}
