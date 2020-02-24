package auth

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

var (
	ecrSvc *ecr.ECR
	region string
)

func init() {
	AuthCmd.AddCommand(EcrCmd)

	EcrCmd.Flags().StringVar(&region, "region", "us-east-1", "aws region to use")

	initSvc()
}

func initSvc() {
	httpClient := &http.Client{}

	cnf := &aws.Config{
		Region:     &region,
		HTTPClient: httpClient,
	}

	sess, err := session.NewSession(cnf)
	if err != nil {
		fmt.Printf("Error creating session: %v", err.Error())
		os.Exit(1)
	}

	ecrSvc = ecr.New(sess)
}

// EcrCmd represents the root of the ECR commands
var EcrCmd = &cobra.Command{
	Use:   "ecr",
	Short: "ECR actions",
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
