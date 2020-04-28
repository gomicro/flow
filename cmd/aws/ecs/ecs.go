package ecs

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

var (
	ecsSvc *ecs.ECS

	region string
)

func init() {
	EcsCmd.Flags().StringVar(&region, "region", "us-east-1", "aws region to use")

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

	ecsSvc = ecs.New(sess)
}

// EcsCmd represents the root of the ecs command
var EcsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "ECS related commands",
}
