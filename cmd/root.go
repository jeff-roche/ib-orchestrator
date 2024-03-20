/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	log "snoman/internal/logger"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// verbose is the optional command that will display INFO logs
var verbose bool

// jsonOutput is the optional command that will display logs as JSON
var jsonOutput bool

// version is an optional command that will display the current release version
var releaseVersion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "snoman",
	Version: releaseVersion,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger = log.Set(verbose, jsonOutput)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	logger.Sync()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Display verbose logs")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "log-json", false, "Format the log output as JSON")

	initCreateCmd()
	initDestroyCmd()
	initGenerateCmd()
}
