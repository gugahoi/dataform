package cmd

import (
	"fmt"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [rds name]",
	Short: "Create a new RDS database",
	Args:  cobra.ExactArgs(1),
	Run:   createFunc,
}

func init() {
	RootCmd.AddCommand(createCmd)
}

func createFunc(cmd *cobra.Command, args []string) {
	session := getAwsSession()
	manager := db.NewManager(rds.New(session))
	instance, err := manager.Create(args[0])
	if err != nil {
		fmt.Printf("Failed to create db: %v", err)
	}
	fmt.Println(instance)
}
