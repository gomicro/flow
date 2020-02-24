package aws

import (
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/cmd/aws/ecr"
)

func init() {
	AwsCmd.AddCommand(ecr.EcrCmd)
}

// AwsCmd represents the root of the aws command
var AwsCmd = &cobra.Command{
	Use:   "aws",
	Short: "AWS related commands",
}
