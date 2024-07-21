package cmd

import (
	"github.com/spf13/cobra"

	"awsexample/cmd/s3"

	"awsexample/cmd/ec2"
)

var Cmd = &cobra.Command{
	Use:   "cmd",
	Short: "cmd ...",

	Run: func(cmd *cobra.Command, args []string) {

		// TODO: Implement command
	},
}

var CmdProfile string

func init() {
	// Flags

	// profile :

	Cmd.PersistentFlags().StringVarP(&CmdProfile, "profile", "", "", "desc")

	Cmd.AddCommand(s3.S3Cmd)

	Cmd.AddCommand(ec2.Ec2Cmd)

}
