/*
Copyright © 2025 Talos-hub wishpersofthenine@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/Talos-hub/ZibraGo/internal/configuration"
	"github.com/spf13/cobra"
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Change settings",
	Long:  `Change path for zip files`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter path for zip files: ")

		var path string
		_, err := fmt.Scan(&path)
		if err != nil {
			fmt.Println("Error to set path")

		}

		err = configuration.NewPath(path)
		if err != nil {
			fmt.Printf("Error to set path: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}
