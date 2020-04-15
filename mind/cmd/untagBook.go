/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// untagBookCmd represents the untagBook command
var untagBookCmd = &cobra.Command{
	Use:        "book",
	Short:      "Remove tag from book",
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"id", "tag"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.UntagBook(args[0], args[1])
	},
}

func init() {
	untagCmd.AddCommand(untagBookCmd)
}
