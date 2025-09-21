/*
Copyright © 2025 Talos-hub wishpersofthenine@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "show current settings",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("settings called")
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	settingsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
