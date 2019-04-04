package cmd

import (
	"filemux/pflag"
	"fmt"
	"os"

	logger "github.com/apsdehal/go-logger"
	"github.com/spf13/cobra"
)

// AppName is the CLI app name.
const AppName = "filemux"

// RootCmd is meant to reused across cmd/*/*.go
var (
	log *logger.Logger

	confFilePath, serverListener string
	logLevel                     = pflag.LogLevel(logger.NoticeLevel)

	rootCmd = &cobra.Command{
		Use: AppName,
	}
)

// Execute execute `rootCmd`.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
