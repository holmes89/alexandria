/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	displayName string
	docType     string
	description string
)

// updateDocumentCmd represents the updateDocument command
var updateDocumentCmd = &cobra.Command{
	Use:        "document",
	Short:      "Update fields of document",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"id"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.UpdateDocument(args[0], displayName, description, docType)
	},
}

func init() {
	updateCmd.AddCommand(updateDocumentCmd)

	updateDocumentCmd.Flags().StringVar(&description, "description", "", "description of doc")
	updateDocumentCmd.Flags().StringVar(&displayName, "display-name", "", "title of doc")
	updateDocumentCmd.Flags().StringVar(&docType, "type", "", "type of doc")
}
