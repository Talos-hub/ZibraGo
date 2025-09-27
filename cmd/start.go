/*
Copyright © 2025 Talos-hub wishpersofthenine@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Talos-hub/ZibraGo/internal/api"
	"github.com/Talos-hub/ZibraGo/internal/archivers"
	"github.com/Talos-hub/ZibraGo/internal/configuration"
	"github.com/Talos-hub/ZibraGo/internal/loggers"
	"github.com/Talos-hub/ZibraGo/internal/services"
	"github.com/Talos-hub/ZibraGo/internal/walkers"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start running backup",
	Long:  `start running backup`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Expected two parameters: name of zip file and path to folder")
			return
		}
		// name of zip file
		name := args[0]
		// folder that it will uploads
		folder := args[1]

		if len(folder) == 0 {
			fmt.Println("Length of path to folder cannot less than ziro")
			return
		}
		if len(name) == 0 {
			fmt.Println("Length of name zip file cannot be")
			return
		}

		// check that name of a zip file is valid
		if !archivers.IsZip(name) {
			fmt.Println("Name of zip file should contatins .zip extenction")
			return
		}

		// setup services...

		// setup logger
		logger, err := loggers.SetupLogger()
		if err != nil {
			fmt.Printf("Error to set up logger: %v\n", err)
			return
		}
		// get config for archiver
		config, err := configuration.GetDir()
		if err != nil {
			if os.IsNotExist(err) {
				logger.Error("To start it requires a configuration file", "error", err)
				fmt.Println("To start it requires a configuration file, you need to use the command: settings")
				return
			}
			logger.Error("Error get configuration file", "error", err)
			fmt.Println("Error get configuration file, you can see a log file with command: log")
			return
		}
		// workers = amount of cpu * 2
		cpu := runtime.NumCPU() * 2
		// setup zip archiver
		zip := archivers.NewZipArchiver(config.ZipDir, name, cpu)
		// setup walker
		wakler := walkers.NewWalker(folder)

		// load credentials
		token, err := api.LoadCredentials()
		if err != nil {
			logger.Error("Error get credentials token", "error", err)
			fmt.Println("Error get credentiald file, you can see a log file with command: log")
			return
		}
		// setup cloud api
		cloud, err := api.NewGoogleCloudApi(token)
		if err != nil {
			logger.Error("Error init cloud service", "error", err)
			fmt.Println("Error init cloud service, you can see a log file with command: log")
			return
		}

		// setup zibra service
		zibra, err := services.NewZibra(wakler, zip, cloud, logger)
		if err != nil {
			logger.Error("Error init zibra service", "error", err)
			fmt.Println("Error init zibra service, you can see a log file with command: log")
			return
		}

		err = zibra.Run()
		if err != nil {
			logger.Error("Error run zibra service", "error", err)
			fmt.Println("Error run zibra service,  you can see a log file with command: log")
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

}
