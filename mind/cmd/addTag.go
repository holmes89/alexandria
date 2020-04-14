/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// addTagCmd represents the addTag command
var addTagCmd = &cobra.Command{
	Use:        "tag",
	Short:      "Create a new tag",
	Long:       `Add a new tag to the system providing a name`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"name"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.CreateTag(args[0])
	},
}

func init() {
	addCmd.AddCommand(addTagCmd)
}
