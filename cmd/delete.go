package cmd

import (
	"fmt"
	"strings"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

var (
	deleteWait bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [rds name]",
	Short: "Delete an existing RDS database",
	Args:  cobra.ExactArgs(1),
	Run:   deleteFunc,
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteWait, "wait", "w", false, "wait for deletion to complete")
	RootCmd.AddCommand(deleteCmd)
}

func deleteFunc(cmd *cobra.Command, args []string) {
	session := getAwsSession()
	manager := db.NewManager(rds.New(session))
	name := args[0]

	fmt.Printf("deleting instance %s\n", name)
	instance, err := manager.Delete(name)
	if err != nil {
		fmt.Printf("failed to delete RDS instance: %v", getAwsError(err))
		return
	}

	if deleteWait {
		go manager.SigHandler()
		state := manager.WaitForFinalState(*instance.Name, 20, 1800)
		for poll := range state {
			if poll.Err != nil {
				if !strings.Contains(poll.Err.Error(), "DBInstanceNotFound") {
					fmt.Printf("error: %v", poll.Err)
					return
				}
			}
			fmt.Printf("%s instance %s\n", poll.Status, name)
		}
	}
}
