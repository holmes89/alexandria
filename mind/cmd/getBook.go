/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// getBookCmd represents the getBook command
var getBookCmd = &cobra.Command{
	Use:   "book",
	Short: "Fetch list of books or single book information",
	Long: `List out all books in library or details on a specific book. List will give you high level information
			like name and upload date.

			Detailed information will reflect all known information about the book. This will need to be done by ID.`,
	Args:       cobra.MaximumNArgs(1),
	ArgAliases: []string{"id"},
	Aliases:    []string{"books"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			results, err := app.FindBooks()
			if err != nil {
				if debug {
					errString := fmt.Errorf("error: %w", err)
					fmt.Fprintln(out, errString.Error())
				}
				return errors.New("unable to fetch books")
			}
			tw := getTabWriter()
			fmt.Fprintf(tw, "\n %s\t%s\t", "ID", "NAME")
			for _, r := range results {
				fmt.Fprintf(tw, "\n %s\t%s\t", r.ID, r.DisplayName)
			}
			fmt.Fprintf(tw, "\n\n")
			tw.Flush()
		} else {
			results, err := app.FindBookByID(args[0])
			if err != nil {
				if debug {
					errString := fmt.Errorf("error: %w", err)
					fmt.Fprintln(out, errString.Error())
				}
				return errors.New("unable to fetch book")
			}
			if results == nil {
				return errors.New("book does not exist")
			}
			b, _ := yaml.Marshal(results)
			fmt.Fprintln(out, string(b))
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getBookCmd)
}
