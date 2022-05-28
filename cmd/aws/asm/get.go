package asm

import (
	"os"

	aws_asm "github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

var ()

func init() {
	AsmCmd.AddCommand(GetCmd)
}

// GetCmd represents the get secret command
var GetCmd = &cobra.Command{
	Use:   "get <secret_name>",
	Short: "Get secret",
	Long:  `Get secret from the Secrets Manager`,
	Args:  cobra.ExactArgs(1),
	Run:   getFunc,
}

func getFunc(cmd *cobra.Command, args []string) {
	name := args[0]

	input := &aws_asm.GetSecretValueInput{
		SecretId: &name,
	}

	resp, err := asmSvc.GetSecretValue(input)
	if err != nil {
		fmt.Printf("Error getting secret: %s", err)
		os.Exit(1)
	}

	fmt.Printf("%v", *resp.SecretString)
}
