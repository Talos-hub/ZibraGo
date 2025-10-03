/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Talos-hub/ZibraGo/internal/configuration"
	"github.com/spf13/cobra"
)

var (
	remove bool
)
var ExCmd = &cobra.Command{
	Use:   "extensions",
	Short: "Manage file extensions",
	Long:  `Add or remove extensions for files that you don't need.`,
	Run: func(cmd *cobra.Command, args []string) {

		if remove {
			// Remove the extensions file
			if err := configuration.RemoveExtensionsFile(); err != nil {
				fmt.Println("Error removing extensions file:", err)
				return
			}
			fmt.Println("Extensions file removed successfully")
			return
		}

		if err := configuration.CreateExtentions(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("\tSuccessful.")
	},
}

func init() {
	rootCmd.AddCommand(ExCmd)
	ExCmd.Flags().BoolVarP(&remove, "rm", "r", false, "Remove the extensions file")
}
