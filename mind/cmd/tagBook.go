/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// tagBookCmd represents the tagBook command
var tagBookCmd = &cobra.Command{
	Use:        "book",
	Short:      "Add a tag to a book",
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"id", "tag"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.TagBook(args[0], args[1])
	},
}

func init() {
	tagCmd.AddCommand(tagBookCmd)
}
