package asm

import (
	"os"
	"sort"

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
	Use:               "get <secret_name>",
	Short:             "Get secret",
	Long:              `Get secret from the Secrets Manager`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getCmdValidArgsFunc,
	Run:               getFunc,
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

func getCmdValidArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	configClient(cmd, args)

	valid := []string{}

	err := asmSvc.ListSecretsPages(&aws_asm.ListSecretsInput{},
		func(page *aws_asm.ListSecretsOutput, lastPage bool) bool {
			for _, sec := range page.SecretList {
				if include(sec.Name) {
					valid = append(valid, *sec.Name)
				}
			}

			return page.NextToken != nil
		})
	if err != nil {
		fmt.Printf("Error listing secrets: %s", err)
		os.Exit(1)
	}

	sort.Strings(valid)

	return valid, cobra.ShellCompDirectiveNoFileComp
}
