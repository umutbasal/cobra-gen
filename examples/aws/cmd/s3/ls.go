package s3

import (
	"fmt"

	"github.com/spf13/cobra"
)

var LsCmd = &cobra.Command{
	Use:   "ls [mybucket]",
	Short: "ls ...",

	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),

	Run: func(cmd *cobra.Command, args []string) {

		mybucket := args[0]
		fmt.Printf("mybucket: %v\n", mybucket)

		// TODO: Implement command
	},
}

var LsPageSize string

func init() {
	// Flags

	// page-size :
	LsCmd.Flags().StringVarP(&LsPageSize, "page-size", "", "", "desc")

	S3Cmd.AddCommand(LsCmd)

}
