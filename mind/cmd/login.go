/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	username string
	password string
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into service",
	Long: `Provide the ability to get login token to use service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.GetAuthToken(username, password)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "username for login")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "password for login")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}
