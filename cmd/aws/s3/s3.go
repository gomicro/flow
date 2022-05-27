package s3

import (
	"os"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/client/aws/session"
	"github.com/gomicro/flow/fmt"
)

var (
	s3Uploader *s3manager.Uploader
)

func init() {
}

// S3Cmd represents the root of the s3 command
var S3Cmd = &cobra.Command{
	Use:              "s3",
	Short:            "S3 related commands",
	PersistentPreRun: configClient,
}

func configClient(cmd *cobra.Command, args []string) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Printf("Error creating session: %v", err.Error())
		os.Exit(1)
	}

	s3Uploader = s3manager.NewUploader(sess)
}
