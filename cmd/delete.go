package cmd

import (
	"fmt"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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
	manager := db.NewManager(rds.New(session.New(&aws.Config{
		Region: aws.String(awsRegion),
	})))

	instance, err := manager.Delete(args[0])
	if err != nil {
		fmt.Printf("Failed to create RDS Instance: %v", err)
	}

	fmt.Println(instance)
}
