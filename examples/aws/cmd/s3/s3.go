package s3

import (
	"github.com/spf13/cobra"
)

var S3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "s3 ...",

	Run: func(cmd *cobra.Command, args []string) {

		// TODO: Implement command
	},
}

func init() {
	// Flags

}
