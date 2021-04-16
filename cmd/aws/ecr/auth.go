package ecr

import (
	"encoding/base64"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

var (
	registries []string
)

func init() {
	EcrCmd.AddCommand(AuthCmd)

	AuthCmd.Flags().StringSliceVar(&registries, "registryID", nil, "aws registry ID to auth with, use flag multiple times to auth with multiple registries")
}

// AuthCmd represents the ECR auth commands
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with ECR",
	Long:  `Authenticate with the AWS ECR service and login to the ECR.`,
	Run:   authFunc,
}

func authFunc(cmd *cobra.Command, args []string) {
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
			cmd := exec.Command("docker", "login", "-u", parts[0], "-p", parts[1], *auth.ProxyEndpoint)
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error executing docker login: %v", auth.String())
			}
		}
	}
}
