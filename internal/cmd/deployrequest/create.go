package deployrequest

import (
	"fmt"

	"github.com/planetscale/cli/internal/cmdutil"
	"github.com/planetscale/cli/internal/printer"
	"github.com/planetscale/planetscale-go/planetscale"

	"github.com/spf13/cobra"
)

// CreateCmd is the command for creating deploy requests.
func CreateCmd(ch *cmdutil.Helper) *cobra.Command {
	var flags struct {
		deployTo string
	}

	cmd := &cobra.Command{
		Use:   "create <database> <branch> [flags]",
		Short: "Create a deploy request from a branch",
		Args:  cmdutil.RequiredArgs("database", "branch"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			database := args[0]
			branch := args[1]

			client, err := ch.Client()
			if err != nil {
				return err
			}

			end := ch.Printer.PrintProgress(fmt.Sprintf("Request deploying of %s branch in %s...", printer.BoldBlue(branch), printer.BoldBlue(database)))
			defer end()

			dr, err := client.DeployRequests.Create(ctx, &planetscale.CreateDeployRequestRequest{
				Organization: ch.Config.Organization,
				Database:     database,
				Branch:       branch,
				IntoBranch:   flags.deployTo,
			})
			if err != nil {
				switch cmdutil.ErrCode(err) {
				case planetscale.ErrNotFound:
					return fmt.Errorf("database %s does not exist in %s",
						printer.BoldBlue(database), printer.BoldBlue(ch.Config.Organization))
				default:
					return cmdutil.HandleError(err)
				}
			}
			end()

			if ch.Printer.Format() == printer.Human {
				number := fmt.Sprintf("#%d", dr.Number)
				ch.Printer.Printf("Deploy request %s successfully created.\n\nView this deploy request in the browser: %s\n", printer.BoldBlue(number), printer.BoldBlue(dr.HtmlURL))
				return nil
			}

			return ch.Printer.PrintResource(toDeployRequest(dr))
		},
	}

	cmd.PersistentFlags().StringVar(&flags.deployTo, "deploy-to", "main", "Branch to deploy the branch. By default it's set to 'main'")

	return cmd
}
