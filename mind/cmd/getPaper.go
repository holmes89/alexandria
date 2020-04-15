/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// getPaperCmd represents the getPaper command
var getPaperCmd = &cobra.Command{
	Use:   "paper",
	Short: "Fetch list of papers or single paper information",
	Long: `List out all papers in library or details on a specific paper. List will give you high level information
			like name and upload date.

			Detailed information will reflect all known information about the paper. This will need to be done by ID.`,
	Args:       cobra.MaximumNArgs(1),
	ArgAliases: []string{"id"},
	Aliases:    []string{"papers"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			results, err := app.FindPapers()
			if err != nil {
				if debug {
					errString := fmt.Errorf("error: %w", err)
					fmt.Fprintln(out, errString.Error())
				}
				return errors.New("unable to fetch papers")
			}
			tw := getTabWriter()
			fmt.Fprintf(tw, "\n %s\t%s\t", "ID", "NAME")
			for _, r := range results {
				fmt.Fprintf(tw, "\n %s\t%s\t", r.ID, r.DisplayName)
			}
			fmt.Fprintf(tw, "\n\n")
			tw.Flush()
		} else {
			results, err := app.FindPaperByID(args[0])
			if err != nil {
				if debug {
					errString := fmt.Errorf("error: %w", err)
					fmt.Fprintln(out, errString.Error())
				}
				return errors.New("unable to fetch paper")
			}
			if results == nil {
				return errors.New("paper does not exist")
			}
			b, _ := yaml.Marshal(results)
			fmt.Fprintln(out, string(b))
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getPaperCmd)
}
