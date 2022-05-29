package asm

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	aws_asm "github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

func init() {
	AsmCmd.AddCommand(SetCmd)
}

// SetCmd represents the set secret command
var SetCmd = &cobra.Command{
	Use:   "set <name> <value>",
	Short: "Set secret",
	Long:  `Set secret in the Secrets Manager. If the secret does not exist it will be created.`,
	Args:  cobra.ExactArgs(2),
	Run:   setFunc,
}

func setFunc(cmd *cobra.Command, args []string) {
	name := args[0]
	value := args[1]

	input := &aws_asm.GetSecretValueInput{
		SecretId: &name,
	}

	exists := true
	_, err := asmSvc.GetSecretValue(input)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			if aerr.Code() != secretsmanager.ErrCodeResourceNotFoundException {
				fmt.Printf("AWS Error: %s", aerr)
				os.Exit(1)
			} else {
				exists = false
			}
		} else {
			fmt.Printf("Error checking for secret's existence: %s", err)
			os.Exit(1)
		}
	}

	if exists {
		input := &secretsmanager.PutSecretValueInput{
			SecretId:     &name,
			SecretString: &value,
		}

		_, err := asmSvc.PutSecretValue(input)
		if err != nil {
			fmt.Printf("Error updating secret: %s", err)
		}
	} else {
		input := &secretsmanager.CreateSecretInput{
			Name:         &name,
			SecretString: &value,
		}

		_, err := asmSvc.CreateSecret(input)
		if err != nil {
			fmt.Printf("Error creating secret: %s", err)
		}
	}
}
