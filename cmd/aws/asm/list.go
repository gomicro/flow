package asm

import (
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

func init() {
	AsmCmd.AddCommand(ListCmd)
}

// ListCmd represents the list secrets command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets",
	Long:  `List secrets from the Secrets Manager`,
	Run:   listFunc,
}

func listFunc(cmd *cobra.Command, args []string) {
	input := &secretsmanager.ListSecretsInput{}

	names := []string{}

	err := asmSvc.ListSecretsPages(input,
		func(page *secretsmanager.ListSecretsOutput, lastPage bool) bool {
			for _, sec := range page.SecretList {
				names = append(names, *sec.Name)
			}

			return page.NextToken != nil
		})
	if err != nil {
		fmt.Printf("Error listing secrets: %s", err)
		os.Exit(1)
	}

	sort.Strings(names)

	for _, n := range names {
		fmt.Printf("%v", n)
	}
}
