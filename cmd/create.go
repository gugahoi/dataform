package cmd

import (
	"fmt"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

var (
	dbMasterUsername     string
	dbMasterPassword     string
	dbInstanceClass      string
	dbEngine             string
	dbEngineVersion      string
	dbMultiAZ            bool
	dbPort               int64
	dbStorageAllocatedGB int64
	dbStorageEncrypted   bool
	dbStorageType        string
	dbStorageIops        int64
	dbSubnetGroup        string
	dbSecurityGroup      string
	dbSecurityGroups     []*string
	dbBackupWindow       string
	dbMaintenanceWindow  string
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
	createCmd.Flags().BoolVarP(&dbStorageEncrypted, "encrypted", "e", true, "db instance encryption")
	createCmd.Flags().StringVarP(&dbStorageType, "type", "t", "gp2", "db storage type")
	createCmd.Flags().Int64VarP(&dbStorageIops, "iops", "i", 0, "db requested iops")
	createCmd.Flags().StringVarP(&dbSubnetGroup, "subnetgroup", "N", "", "db subnet group name")
	createCmd.Flags().StringVarP(&dbSecurityGroup, "securitygroup", "S", "", "db security group id")
	createCmd.Flags().StringVarP(&dbBackupWindow, "backup", "B", "", "db preferred backup window")
	createCmd.Flags().StringVarP(&dbMaintenanceWindow, "maintenance", "M", "", "db preferred maintenance window")
	RootCmd.AddCommand(createCmd)
}

func createFunc(cmd *cobra.Command, args []string) {
	session := getAwsSession()
	manager := db.NewManager(rds.New(session))
	name := args[0]

	dbinput := &db.DB{
		MasterUsername:     &dbMasterUsername,
		MasterUserPassword: &dbMasterPassword,
		Engine:             &dbEngine,
		EngineVersion:      &dbEngineVersion,
		DBInstanceClass:    &dbInstanceClass,
		Name:               &name,
		MultiAZ:            &dbMultiAZ,
		Port:               &dbPort,
		StorageAllocatedGB: &dbStorageAllocatedGB,
	}
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

	instance, err := manager.Create(dbinput)
	if err != nil {
		fmt.Printf("Failed to create RDS Instance: %v", getAwsError(err))
		return
	}

	fmt.Printf("%s\t%s\t%s\n", *instance.Name, *instance.ARN, *instance.Status)
}
