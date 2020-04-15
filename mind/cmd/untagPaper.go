/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// untagPaperCmd represents the untagPaper command
var untagPaperCmd = &cobra.Command{
	Use:        "paper",
	Short:      "Remove tag from paper",
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"id", "tag"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.UntagPaper(args[0], args[1])
	},
}

func init() {
	untagCmd.AddCommand(untagPaperCmd)
}
