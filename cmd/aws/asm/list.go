package asm

import (
	"os"
	"sort"
	"strings"

	aws_asm "github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

var (
	filters []string
)

func init() {
	AsmCmd.AddCommand(ListCmd)

	ListCmd.Flags().StringSliceVarP(&filters, "filter", "f", []string{}, "filter the list by the specified value")
}

// ListCmd represents the list secrets command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets",
	Long:  `List secrets from the Secrets Manager`,
	Run:   listFunc,
}

func listFunc(cmd *cobra.Command, args []string) {
	input := &aws_asm.ListSecretsInput{}

	names := []string{}

	err := asmSvc.ListSecretsPages(input,
		func(page *aws_asm.ListSecretsOutput, lastPage bool) bool {
			for _, sec := range page.SecretList {
				if include(sec.Name) {
					names = append(names, *sec.Name)
				}
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

func include(secret *string) bool {
	for _, filter := range filters {
		if !strings.Contains(strings.ToLower(*secret), filter) {
			return false
		}
	}

	return true
}
