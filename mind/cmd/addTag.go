/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addTagCmd represents the addTag command
var addTagCmd = &cobra.Command{
	Use:        "addTag",
	Short:      "Create a new tag",
	Long:       `Add a new tag to the system providing a name`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"name"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addTag called")
	},
}

func init() {
	addCmd.AddCommand(addTagCmd)
}
