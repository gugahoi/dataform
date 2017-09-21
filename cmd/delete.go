package cmd

import (
	"fmt"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [rds name]",
	Short: "Delete an existing RDS database",
	Args:  cobra.ExactArgs(1),
	Run:   deleteFunc,
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}

func deleteFunc(cmd *cobra.Command, args []string) {
	session := getAwsSession()
	manager := db.NewManager(rds.New(session))
	instance, err := manager.Delete(args[0])
	if err != nil {
		fmt.Printf("Failed to delete RDS Instance: %s", getAwsError(err))
	}

	fmt.Println(instance)
}
