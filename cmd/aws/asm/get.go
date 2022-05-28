package asm

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	aws_asm "github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

func init() {
	AsmCmd.AddCommand(GetCmd)
}

// GetCmd represents the get secret command
var GetCmd = &cobra.Command{
	Use:               "get <secret_name>",
	Short:             "Get secret",
	Long:              `Get secret from the Secrets Manager. Uses of periods in the secret name will be used to denote selecting into a json object if it is returned.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getCmdValidArgsFunc,
	Run:               getFunc,
}

func getFunc(cmd *cobra.Command, args []string) {
	name := args[0]

	parts := strings.Split(name, ".")

	input := &aws_asm.GetSecretValueInput{}

	if len(parts) > 1 {
		input.SecretId = &parts[0]
	} else {
		input.SecretId = &name
	}

	resp, err := asmSvc.GetSecretValue(input)
	if err != nil {
		fmt.Printf("Error getting secret: %s", err)
		os.Exit(1)
	}

	if len(parts) > 1 {
		var heap map[string]interface{}
		err = json.Unmarshal([]byte(*resp.SecretString), &heap)
		if err != nil {
			fmt.Printf("Error unmarshalling secret json: %v", err.Error())
			os.Exit(1)
		}

		val, err := getJsonValue(heap, parts[1:])
		if err != nil {
			fmt.Printf("Error getting json value: %v", err.Error())
			os.Exit(1)
		}

		fmt.Printf("%v", val)

	} else {
		fmt.Printf("%v", *resp.SecretString)
	}
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
