package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd ...
var RootCmd = &cobra.Command{
	Use:   "dfm",
	Short: "Give me ALL your RDS",
}

var (
	awsRegion string
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&awsRegion, "region", "", "", "AWS Region")
}

// Execute ...
func Execute() {
	RootCmd.Execute()
}
