package ecr

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

var (
	region string
)

func init() {
	EcrCmd.AddCommand(AuthCmd)

	AuthCmd.Flags().StringVar(&region, "region", "us-east-1", "aws region to use")
}

// AuthCmd represents the root of the ECR commands
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with ECR",
	Long:  `Authenticate with the AWS ECR service`,
	Run:   ecrFunc,
}

func ecrFunc(cmd *cobra.Command, args []string) {
	input := &ecr.GetAuthorizationTokenInput{
		//RegistryIds: []*string{},
	}

	auths, err := ecrSvc.GetAuthorizationToken(input)
	if err != nil {
		fmt.Printf("Error getting ecr auth token: %v", err.Error())
		os.Exit(1)
	}

	if auths == nil {
		fmt.Printf("Empty reponse from ecr auth")
		os.Exit(1)
	}

	for _, auth := range auths.AuthorizationData {
		tkn, _ := base64.StdEncoding.DecodeString(*auth.AuthorizationToken)

		parts := strings.SplitN(string(tkn), ":", 2)
		fmt.Printf("docker login -u %v -p %v %v", parts[0], parts[1], *auth.ProxyEndpoint)
	}
}