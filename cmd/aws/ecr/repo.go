package ecr

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

func init() {
	EcrCmd.AddCommand(RepoCmd)
}

// RepoCmd represents the repo command
var RepoCmd = &cobra.Command{
	Use:   "repo <repository>",
	Short: "Create ECR repo",
	Long:  `Create an ECR repo if it does not exist.`,
	Args:  cobra.ExactArgs(1),
	Run:   repoFunc,
}

func repoFunc(cmd *cobra.Command, args []string) {
	repo := args[0]
	if !strings.Contains(repo, ".dkr.ecr.") {
		fmt.Printf("Repo must be fully qualified ECR repo")
		os.Exit(1)
	}

	parts := strings.Split(repo, ".")
	if len(parts) != 6 {
		fmt.Printf("Repo missing information")
		os.Exit(1)
	}

	regID := parts[0]

	image := strings.TrimPrefix(parts[len(parts)-1], "com/")
	repoParts := strings.Split(image, ":")
	repoName := repoParts[0]

	if repoExists(regID, repoName) {
		fmt.Verbosef("Repo found, nothing to do")
		return
	}

	createRepo(repoName)
}

func repoExists(regID, repoName string) bool {
	input := &ecr.DescribeRepositoriesInput{
		RegistryId:      &regID,
		RepositoryNames: []*string{&repoName},
	}

	search, err := ecrSvc.DescribeRepositories(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != ecr.ErrCodeRepositoryNotFoundException {
				fmt.Printf("Error getting ecr repos: %v", err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error getting ecr repos: %v", err.Error())
			os.Exit(1)
		}
	}

	if len(search.Repositories) == 0 {
		return false
	}

	return true
}

func createRepo(repoName string) {
	input := &ecr.CreateRepositoryInput{
		RepositoryName: &repoName,
	}

	_, err := ecrSvc.CreateRepository(input)
	if err != nil {
		fmt.Printf("Failed to create repository: %v", err.Error())
		os.Exit(1)
	}
}
