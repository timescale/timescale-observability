package grafana

import (
	"fmt"

	"github.com/spf13/cobra"
	root "github.com/timescale/tobs/cli/cmd"
)

// grafanaGetPasswordCmd represents the grafana get-password command
var grafanaGetPasswordCmd = &cobra.Command{
	Use:   "get-password",
	Short: "Gets the admin password for Grafana",
	Args:  cobra.ExactArgs(0),
	RunE:  grafanaGetPassword,
}

func init() {
	grafanaCmd.AddCommand(grafanaGetPasswordCmd)
}

func grafanaGetPassword(cmd *cobra.Command, args []string) error {
	var err error

	secret, err := kubeClient.KubeGetSecret(root.Namespace, root.HelmReleaseName+"-grafana")
	if err != nil {
		return fmt.Errorf("could not get Grafana password: %w", err)
	}

	pass := secret.Data["admin-password"]
	fmt.Println(string(pass))

	return nil
}
