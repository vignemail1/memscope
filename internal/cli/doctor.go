package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose memscope capabilities and system compatibility",
		Long:  "Checks system capabilities and reports on memscope's ability to access hardware",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "System Diagnostics (placeholder)")
			return nil
		},
	}
}