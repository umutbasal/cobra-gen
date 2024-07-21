package ec2

import (
	"github.com/spf13/cobra"
)

var Ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "ec2 ...",

	Run: func(cmd *cobra.Command, args []string) {

		// TODO: Implement command
	},
}

func init() {
	// Flags

}
