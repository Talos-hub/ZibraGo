/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Talos-hub/ZibraGo/internal/configuration"
	"github.com/spf13/cobra"
)

var ExCmd = &cobra.Command{
	Use:   "extentions",
	Short: "Add extentions",
	Long:  `Add extentions for files that you don't need.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := configuration.CreateExtentions(); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(ExCmd)
}
