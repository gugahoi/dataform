package cmd

import (
	"fmt"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var listCmd = &cobra.Command{
	Use:   "list [rds name]",
	Short: "list RDS databases",
	Run:   listFunc,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

// listFunc return the name and status of current rds instances
func listFunc(cmd *cobra.Command, args []string) {
	session := getAwsSession()
	manager := db.NewManager(rds.New(session))
	if len(args) > 0 {
		for _, name := range args {
			list, err := manager.Stat(name)
			if err != nil {
				fmt.Printf("%s: %s\n", name, getAwsError(err))
			} else {
				printInstances(list)
			}
		}
	} else {
		list, err := manager.StatAll()
		if err != nil {
			fmt.Printf("error: failed to list dbs: %v\n", getAwsError(err))
			return
		} else {
			printInstances(list)
		}
	}
}

// printInstances short form output of the rds instances
func printInstances(list *rds.DescribeDBInstancesOutput) {
	for _, db := range list.DBInstances {
		fmt.Printf("%s: %s\n", *db.DBInstanceIdentifier, *db.DBInstanceStatus)
	}
}
