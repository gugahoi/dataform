package cmd

import (
	"fmt"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all RDS instances in a region",
	Run:   listFunc,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

// listFunc return the name and status of current rds instances
func listFunc(cmd *cobra.Command, args []string) {
	session := getAwsSession()
	manager := db.NewManager(rds.New(session))

	results, err := manager.List()
	if err != nil {
		fmt.Printf("Failed listing RDS instances: %s\n", getAwsError(err))
		return
	}

	for _, db := range results {
		fmt.Printf("%s\t%s\t%s\n", *db.Name, *db.Status, *db.ARN)
	}
}
