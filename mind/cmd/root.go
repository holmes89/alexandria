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
	"fmt"
	"github.com/Holmes89/alexandria/mind/internal"
	"github.com/spf13/cobra"
	"io"
	"os"
	"text/tabwriter"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	app        *internal.App
	out        io.Writer
	uploadPath string
	name       string
	debug      bool
)

func getTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(out, 0, 0, 4, ' ', 0)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mind",
	Short: "Tool to support my brain",
	Long: `I have lots of elements of things I want to keep track of map and trace. 
This tool should help me achieve some level of organization to support the growth of ideas and planning.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(out, err)
		os.Exit(1)
	}
}

func init() {
	// init config and writer, could support file io in the future
	initWriter()
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mind.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "verbose", false, "extra debug information used")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func initWriter() {
	out = os.Stdout
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(out, err)
			os.Exit(1)
		}

		configFile := fmt.Sprintf("%s/.mind.yml", home)
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			file, err := os.Create(configFile)
			if err != nil {
				fmt.Fprintln(out, "unable to create config file")
				os.Exit(1)
			}
			fmt.Println(file.Name())
			file.Close()

		}

		// Search config in home directory with name ".mind" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mind")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Fprintln(out, "Using config file:", viper.ConfigFileUsed())
		}

	}

	app = &internal.App{}
	if err := viper.Unmarshal(app); err != nil {
		fmt.Fprintln(out, "unable to load config")
		os.Exit(1)
	}

	app.Config = viper.GetViper()

}
