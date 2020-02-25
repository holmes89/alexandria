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

// setConfigCmd represents the setConfig command
var setConfigCmd = &cobra.Command{
	Use:        "set",
	Short:      "Set configuration value",
	Long:       `Allow for custom configuration of the application. Currently only ability to set endpoint`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"key", "value"},
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]
		if err := app.SetConfigurationValue(key, value); err != nil {
			return errors.New("unable to set key")
		}
		fmt.Fprintf(out, "set value for %s\n", key)
		return nil
	},
}

func init() {
	configCmd.AddCommand(setConfigCmd)
}
