/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// deleteDocumentCmd represents the deleteDocument command
var deleteDocumentCmd = &cobra.Command{
	Use:        "document",
	Short:      "Remove document from library",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"id"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.DeleteDocument(args[0])
	},
}

func init() {
	deleteCmd.AddCommand(deleteDocumentCmd)

}
