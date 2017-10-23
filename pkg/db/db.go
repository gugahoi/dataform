package db

import (
	"fmt"
	"math/rand"
	time "time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

var (
	// Defaults are rds instance defaults
	Defaults = &DB{
		DBInstanceClass:    aws.String("db.t2.small"),
		CopyTagsToSnapshot: aws.Bool(true),
		Engine:             aws.String("postgres"),
		EngineVersion:      aws.String("9.6.3"),
		MultiAZ:            aws.Bool(false),
		Port:               aws.Int64(5432),
		StorageAllocatedGB: aws.Int64(5),
		StorageEncrypted:   aws.Bool(true),
		StorageType:        aws.String("gp2"),
	}

	errInvalidUsernamePassword     = fmt.Errorf("username and password cannot be empty")
	errDbNameMissing               = fmt.Errorf("error: required DB field Name is missing")
	errDbMasterUsernameMissing     = fmt.Errorf("error: required DB field MasterUsername is missing")
	errDbMasterUserPasswordMissing = fmt.Errorf("error: required DB field MasterUserPassword is missing")
)

// Manager uses a svc to talk to AWS RDS
type Manager struct {
	Client rdsiface.RDSAPI
}

// NewManager returns a pointer to a Manager struct.
// The supplied svc is used to make calls to AWS RDS Service.
func NewManager(svc rdsiface.RDSAPI) *Manager {
	return &Manager{svc}
}

// Create an RDS Instance from a supplied DB object
func (r *Manager) Create(db *DB) (*DB, error) {
	database, err := validateDBInput(db)
	if err != nil {
		return nil, err
	}
	dbInput := &rds.CreateDBInstanceInput{
		AllocatedStorage:     database.StorageAllocatedGB,
		DBInstanceClass:      database.DBInstanceClass,
		DBInstanceIdentifier: database.Name,
		Engine:               database.Engine,
		EngineVersion:        database.EngineVersion,
		MasterUserPassword:   database.MasterUserPassword,
		MasterUsername:       database.MasterUsername,
		MultiAZ:              database.MultiAZ,
		Port:                 database.Port,
		DBSubnetGroupName:    database.SubnetGroupName,
		VpcSecurityGroupIds:  database.SecurityGroups,
		StorageEncrypted:     database.StorageEncrypted,
		StorageType:          database.StorageType,
		Iops:                 database.StorageIops,
	}

	if database.KMSKeyArn != nil {
		dbInput.KmsKeyId = database.KMSKeyArn
	}
	if database.PreferredBackupWindow != nil {
		dbInput.PreferredBackupWindow = database.PreferredBackupWindow
	}
	if database.PreferredMaintenanceWindow != nil {
		dbInput.PreferredMaintenanceWindow = database.PreferredMaintenanceWindow
	}
	if database.Tags != nil {
		tags := make([]*rds.Tag, 0, 10)
		for _, v := range database.Tags {
			tags = append(tags, &rds.Tag{
				Key:   v.Key,
				Value: v.Value,
			})
		}
		dbInput.Tags = tags
	}

	result, err := r.Client.CreateDBInstance(dbInput)
	if err != nil {
		return nil, err
	}

	return FromDBInstance(result.DBInstance), nil
}

func validateDBInput(db *DB) (*DB, error) {
	if db.Name == nil {
		return nil, errDbNameMissing
	}
	if db.MasterUsername == nil {
		return nil, errDbMasterUsernameMissing
	}
	if len(*db.MasterUsername) == 0 {
		return nil, errDbMasterUsernameMissing
	}
	if db.MasterUserPassword == nil {
		return nil, errDbMasterUserPasswordMissing
	}
	if len(*db.MasterUserPassword) == 0 {
		return nil, errDbMasterUserPasswordMissing
	}

	// cannot override this one for now
	db.Engine = Defaults.Engine
	db.CopyTagsToSnapshot = Defaults.CopyTagsToSnapshot

	if db.MultiAZ == nil {
		db.MultiAZ = Defaults.MultiAZ
	}
	if db.DBInstanceClass == nil {
		db.DBInstanceClass = Defaults.DBInstanceClass
	}
	if db.EngineVersion == nil {
		db.EngineVersion = Defaults.EngineVersion
	}
	if db.StorageAllocatedGB == nil {
		db.StorageAllocatedGB = Defaults.StorageAllocatedGB
	}
	if db.StorageType == nil {
		db.StorageType = Defaults.StorageType
	}
	if db.StorageEncrypted == nil {
		db.StorageEncrypted = Defaults.StorageEncrypted
	}

	return db, nil
}

// Delete an RDS Instance with the given name
func (r *Manager) Delete(name string) (*DB, error) {
	dbInstanceInput := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(name),
		SkipFinalSnapshot:    aws.Bool(false),
	}

	result, err := r.Client.DeleteDBInstance(dbInstanceInput)
	if err != nil {
		return nil, err
	}

	return FromDBInstance(result.DBInstance), nil
}

// Stat returns the status of an RDS Instance
func (r *Manager) Stat(name string) (*DB, error) {
	dbInstanceInput := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(name),
	}

	result, err := r.Client.DescribeDBInstances(dbInstanceInput)
	if err != nil {
		return nil, err
	}

	if len(result.DBInstances) == 0 {
		return nil, nil
	}

	return FromDBInstance(result.DBInstances[0]), nil
}

// List returns the status of all RDS Instances
func (r *Manager) List() ([]*DB, error) {
	dbInstanceInput := &rds.DescribeDBInstancesInput{}

	result, err := r.Client.DescribeDBInstances(dbInstanceInput)
	if err != nil {
		return nil, err
	}

	return FromDBInstances(result.DBInstances), nil
}

// generateRandomString receives a size and a string of allowed characters and generates a random string of the given size
func generateRandomString(strlen int, allowedChars string, t clock) string {
	rsource := rand.New(rand.NewSource(t.Now().UnixNano()))
	result := make([]byte, strlen)
	for i := range result {
		result[i] = allowedChars[rsource.Intn(len(allowedChars))]
	}
	return string(result)
}

// GenerateRandomPassword receives a size and generates a random password of that size
func (r *Manager) GenerateRandomPassword(strlen int) string {
	const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFHIJKLMNOPQRSTUVWXYZ0123456789!#$%^&*()-+="
	return generateRandomString(strlen, allowedChars, actualClock{})
}

// GenerateRandomUsername receives a size and generates a random username of that size
func (r *Manager) GenerateRandomUsername(strlen int) string {
	const allowedChars = "abcdefghijklmnopqrstuvwxyz"
	return generateRandomString(strlen, allowedChars, actualClock{})
}

// clock allows us to mock out time.Now in our tests
type clock interface {
	Now() time.Time
}

type actualClock struct{}

func (actualClock) Now() time.Time {
	return time.Now()
}
