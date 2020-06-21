package ecr

import (
	"net/http"
	"os"

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
	EcrCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "aws region to use")
}

// EcrCmd represents the root of the auth command
var EcrCmd = &cobra.Command{
	Use:              "ecr",
	Short:            "ECR related commands",
	PersistentPreRun: configClient,
}

func configClient(cmd *cobra.Command, args []string) {
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
