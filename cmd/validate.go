package cmd

import (
	"limux/processor"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check the provided configuration",
	Run: func(c *cobra.Command, args []string) {
		_, err := initConfig(confFilePath)
		if err != nil {
			if _, ok := err.(processor.EnvironmentError); !ok {
				fmt.Println(err)

				os.Exit(1)
			}
		}

		fmt.Println("The provided configuration is valid.")

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&confFilePath, "conf-file", "c", "", "configuration file path")
	validateCmd.MarkFlagRequired("conf-file")
}
