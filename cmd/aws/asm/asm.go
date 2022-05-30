package asm

import (
	"os"

	aws_asm "github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/client/aws/session"
	"github.com/gomicro/flow/fmt"
)

var (
	asmSvc *aws_asm.SecretsManager
)

func init() {
}

// AsmCmd represents the root of the asm command
var AsmCmd = &cobra.Command{
	Use:              "asm",
	Short:            "asm related commands",
	PersistentPreRun: configClient,
}

func configClient(cmd *cobra.Command, args []string) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Printf("Error creating session: %v", err.Error())
		os.Exit(1)
	}

	asmSvc = aws_asm.New(sess)
}
