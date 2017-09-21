package cmd

import (
	"fmt"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all RDS databases in a region",
	Run:   listFunc,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

// listFunc return the name and status of current rds instances
func listFunc(cmd *cobra.Command, args []string) {
	session := getAwsSession()
	manager := db.NewManager(rds.New(session))

	list, err := manager.List()
	if err != nil {
		fmt.Printf("Failed listing: %s\n", getAwsError(err))
		return
	}

	for idx, db := range list.DBInstances {
		fmt.Printf("%d\t%s\t%s\n", idx+1, *db.DBInstanceIdentifier, *db.DBInstanceStatus)
	}
}
