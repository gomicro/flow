package profile

import (
	gofmt "fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gomicro/flow/client/aws/config"
	"github.com/gomicro/flow/fmt"
)

const (
	awsProfileEnv = "AWS_PROFILE"
)

// ProfileCmd represents the profile command
var ProfileCmd = &cobra.Command{
	Use:               "profile [use_profile]",
	Short:             "see current profile or set profile to use",
	Long:              `See the current profile that is set, or set what profile to use.`,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: profileCmdValidArgsFunc,
	Run:               profileFunc,
}

func profileFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Printf("Current Profile: %s", os.Getenv(awsProfileEnv))
		return
	}

	gofmt.Printf("export %s=%s", awsProfileEnv, args[0])
}

func profileCmdValidArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ps, err := config.GetProfiles()
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	vs := make([]string, len(ps))
	for _, p := range ps {
		v := p.Name
		switch {
		case p.AwsAccountID != "":
			v = gofmt.Sprintf("%s\t%s", p.Name, p.AwsAccountID)
		case p.RoleArn != "":
			id := strings.Split(strings.TrimPrefix(p.RoleArn, "arn:aws:iam::"), ":")[0]
			v = gofmt.Sprintf("%s\t%s", p.Name, id)
		}

		vs = append(vs, v)
	}

	return vs, cobra.ShellCompDirectiveNoFileComp
}
