package cmd

import (
	"fmt"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var statCmd = &cobra.Command{
	Use:   "stat [rds name]",
	Short: "Describe an RDS instance",
	Args:  cobra.ExactArgs(1),
	Run:   statFunc,
}

func init() {
	RootCmd.AddCommand(statCmd)
}

// statFunc return the name and status of current rds instances
func statFunc(cmd *cobra.Command, args []string) {
	session := getAwsSession()
	manager := db.NewManager(rds.New(session))
	name := args[0]

	i, err := manager.Stat(name)
	if err != nil {
		fmt.Printf("%s: %s\n", name, getAwsError(err))
		return
	}

	fmt.Println(i)
}
