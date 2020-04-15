/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// tagPaperCmd represents the tagPaper command
var tagPaperCmd = &cobra.Command{
	Use:        "paper",
	Short:      "Add a tag to a paper",
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"id", "tag"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.TagPaper(args[0], args[1])
	},
}

func init() {
	tagCmd.AddCommand(tagPaperCmd)
}
