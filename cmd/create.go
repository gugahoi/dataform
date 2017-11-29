package cmd

import (
	"fmt"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

var (
	dbMasterUsername      string
	dbMasterPassword      string
	dbInstanceClass       string
	dbEngine              string
	dbEngineVersion       string
	dbMultiAZ             bool
	dbPort                int64
	dbStorageAllocatedGB  int64
	dbStorageEncrypted    bool
	dbStorageType         string
	dbStorageIops         int64
	dbSubnetGroup         string
	dbSecurityGroup       string
	dbSecurityGroups      []*string
	dbBackupWindow        string
	dbMaintenanceWindow   string
	dbBackupRetentionDays int64
	createWait            bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [rds name]",
	Short: "Create a new RDS database",
	Args:  cobra.ExactArgs(1),
	Run:   createFunc,
}

func init() {
	createCmd.Flags().StringVarP(&dbMasterUsername, "username", "u", "admin", "db master username")
	createCmd.Flags().StringVarP(&dbMasterPassword, "password", "p", "", "db master password")
	createCmd.Flags().StringVarP(&dbEngine, "engine", "e", "postgres", "db engine")
	createCmd.Flags().StringVarP(&dbEngineVersion, "version", "v", "9.6.3", "db engine version")
	createCmd.Flags().StringVarP(&dbInstanceClass, "class", "c", "db.t2.small", "db instance class/size")
	createCmd.Flags().BoolVarP(&dbMultiAZ, "multiaz", "m", false, "db instance multiAZ")
	createCmd.Flags().Int64VarP(&dbPort, "port", "P", 5432, "db instance port number")
	createCmd.Flags().Int64VarP(&dbStorageAllocatedGB, "size", "s", 5, "db storage size allocated in GB")
	createCmd.Flags().BoolVarP(&dbStorageEncrypted, "encrypted", "E", true, "db instance encryption")
	createCmd.Flags().StringVarP(&dbStorageType, "type", "t", "gp2", "db storage type")
	createCmd.Flags().Int64VarP(&dbStorageIops, "iops", "i", 0, "db requested iops")
	createCmd.Flags().StringVarP(&dbSubnetGroup, "subnetgroup", "N", "", "db subnet group name")
	createCmd.Flags().StringVarP(&dbSecurityGroup, "securitygroup", "S", "", "db security group id")
	createCmd.Flags().StringVarP(&dbBackupWindow, "backup", "B", "", "db preferred backup window")
	createCmd.Flags().StringVarP(&dbMaintenanceWindow, "maintenance", "M", "", "db preferred maintenance window")
	createCmd.Flags().Int64VarP(&dbBackupRetentionDays, "backupretentiondays", "d", 20, "db preferred maintenance window")
	createCmd.Flags().BoolVarP(&createWait, "wait", "w", false, "wait for creation to complete")
	RootCmd.AddCommand(createCmd)
}

func createFunc(cmd *cobra.Command, args []string) {
	session := getAwsSession()
	manager := db.NewManager(rds.New(session))
	name := args[0]

	dbinput := &db.DB{}
	dbinput.MasterUsername = &dbMasterUsername
	dbinput.MasterUserPassword = &dbMasterPassword
	dbinput.Engine = &dbEngine
	dbinput.EngineVersion = &dbEngineVersion
	dbinput.DBInstanceClass = &dbInstanceClass
	dbinput.Name = &name
	dbinput.Port = &dbPort
	dbinput.StorageAllocatedGB = &dbStorageAllocatedGB
	dbinput.MultiAZ = &dbMultiAZ
	dbinput.BackupRetentionPeriod = &dbBackupRetentionDays

	if dbStorageIops > 0 {
		dbinput.StorageIops = &dbStorageIops
	}
	if dbSubnetGroup != "" {
		dbinput.SubnetGroupName = &dbSubnetGroup
	}
	if dbSecurityGroup != "" {
		dbSecurityGroups = make([]*string, 0, 5)
		dbSecurityGroups = append(dbSecurityGroups, &dbSecurityGroup)
	}
	if dbBackupWindow != "" {
		dbinput.PreferredBackupWindow = &dbBackupWindow
	}
	if dbMaintenanceWindow != "" {
		dbinput.PreferredMaintenanceWindow = &dbMaintenanceWindow
	}

	fmt.Printf("creating instance %s\n", *dbinput.Name)

	instance, err := manager.CreateDBInstance(dbinput, db.Development)

	if err != nil {
		fmt.Printf("failed to create instance: %v", getAwsError(err))
		return
	}
	if createWait {
		go manager.SigHandler()
		status := manager.WaitForFinalState(*instance.Name, 20, 1800)
		for poll := range status {
			if poll.Err != nil {
				fmt.Printf("instance transitioned to error condition: %v", err)
				return
			}
			fmt.Printf("%s instance %s\n", poll.Status, *instance.Name)
		}
	}
	fmt.Printf("created %s %s\n", *instance.Name, *instance.ARN)
}
