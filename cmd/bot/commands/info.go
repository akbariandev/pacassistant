package commands

import (
	"fmt"
	"runtime"

	"github.com/akbariandev/pacassistant/version"

	"github.com/spf13/cobra"
)

var (
	BuildDate string
	CommitID  string
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "bot service information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf(infoTemplate, version.ServiceName, version.ServiceDescription, version.Version, BuildDate, CommitID,
			runtime.Version(), runtime.GOOS, runtime.GOARCH)

		return nil
	},
}

var infoTemplate = `Service Name: %s
Description: %s
Version: %s
Build Date: %s
Commit ID: %s
Go version: %v
OS built app: %v - %v
`
