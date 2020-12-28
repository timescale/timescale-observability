package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// helmShowValuesCmd represents the helm show-values command
var helmShowValuesCmd = &cobra.Command{
	Use:   "show-values",
	Short: "Prints the default Observability Stack values to console",
	Args:  cobra.ExactArgs(0),
	RunE:  helmShowValues,
}

func init() {
	helmCmd.AddCommand(helmShowValuesCmd)
}

func helmShowValues(cmd *cobra.Command, args []string) error {
	var err error

	var showvalues *exec.Cmd
	if DEVEL {
		showvalues = exec.Command("helm", "show", "values", "timescale/tobs", "--devel")
	} else {
		showvalues = exec.Command("helm", "show", "values", "timescale/tobs")
	}
	showvalues.Stderr = os.Stderr

	var out []byte
	out, err = showvalues.Output()
	if err != nil {
		return fmt.Errorf("could not get Helm values: %w", err)
	}

	fmt.Print(string(out))

	return nil
}
