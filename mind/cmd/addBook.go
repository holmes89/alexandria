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

	"github.com/spf13/cobra"
)

// addBookCmd represents the addBook command
var addBookCmd = &cobra.Command{
	Use:   "book",
	Short: "Upload book to the service",
	Long:  `Add book to library from local file system providing a path and a name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := app.UploadBook(uploadPath, name); err != nil {
			if debug {
				fmt.Fprintln(out, err.Error())
			}
			return errors.New("unable to upload file")
		}
		fmt.Fprintln(out, "file successfully uploaded")
		return nil
	},
}

func init() {
	addCmd.AddCommand(addBookCmd)

	addBookCmd.Flags().StringVarP(&uploadPath, "path", "p", "", "filepath to upload")
	addBookCmd.Flags().StringVarP(&name, "name", "n", "", "display name of the file")
	addBookCmd.MarkFlagRequired("path")
	addBookCmd.MarkFlagRequired("name")
}
