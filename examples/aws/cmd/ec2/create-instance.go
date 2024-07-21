package ec2

import (
	"github.com/spf13/cobra"
)

var CreateInstanceCmd = &cobra.Command{
	Use:   "create-instance",
	Short: "create-instance ...",

	Run: func(cmd *cobra.Command, args []string) {

		// TODO: Implement command
	},
}

func init() {
	// Flags

	Ec2Cmd.AddCommand(CreateInstanceCmd)

}
