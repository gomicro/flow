package auth

import (
	"github.com/spf13/cobra"
)

// AuthCmd represents the root of the auth command
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate for running commands",
	Long:  `Authenticate with the cloud provider for a given service.`,
}
