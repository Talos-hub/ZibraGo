/*
Copyright © 2025 Talos-hub wishpersofthenine@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Talos-hub/ZibraGo/internal/loggers"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Shows all logs",
	Long:  `Shows all logs`,
	Run: func(cmd *cobra.Command, args []string) {
		err := loggers.Show()
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Log file not exists")
				return
			}
		}

		fmt.Printf("Error show logs %v\n", err)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
