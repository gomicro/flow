package ecr

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

func init() {
	EcrCmd.AddCommand(GetAuthCmd)

	GetAuthCmd.Flags().StringSliceVar(&registries, "registryID", nil, "aws registry ID to auth with, use flag multiple times to auth with multiple registries")
}

// GetAuthCmd represents the ECR get-auth commands
var GetAuthCmd = &cobra.Command{
	Use:   "get-auth",
	Short: "Get ECR authentication creds",
	Long:  `Get the credentials for authenticating with the AWS ECR service.`,
	Run:   getAuthFunc,
}

func getAuthFunc(cmd *cobra.Command, args []string) {
	rs := make([]*string, len(registries), cap(registries))

	for i := range registries {
		rs[i] = &registries[i]
	}

	input := &ecr.GetAuthorizationTokenInput{}

	if len(registries) > 0 {
		input.RegistryIds = rs
	}

	auths, err := ecrSvc.GetAuthorizationToken(input)
	if err != nil {
		fmt.Printf("Error getting ecr auth token: %v", err.Error())
		os.Exit(1)
	}

	if auths == nil {
		fmt.Printf("Empty reponse from ecr auth")
		os.Exit(1)
	} else {
		for _, auth := range auths.AuthorizationData {
			tkn, _ := base64.StdEncoding.DecodeString(*auth.AuthorizationToken)

			parts := strings.SplitN(string(tkn), ":", 2)
			fmt.Printf("docker login -u %v -p %v %v", parts[0], parts[1], *auth.ProxyEndpoint)
		}
	}
}
