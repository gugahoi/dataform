package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
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

// Create an RDS Instance with the given name.
func (r *Manager) Create(name string) (*DB, error) {
	dbInput := &rds.CreateDBInstanceInput{
		AllocatedStorage:     aws.Int64(5),
		DBInstanceClass:      aws.String("db.t2.micro"),
		DBInstanceIdentifier: aws.String(name),
		Engine:               aws.String("postgres"),
		MasterUserPassword:   aws.String("mypassword"),
		MasterUsername:       aws.String("masteruser"),
	}

	result, err := r.Client.CreateDBInstance(dbInput)
	if err != nil {
		return nil, err
	}

	return FromDBInstance(result.DBInstance), nil
}

// Delete an RDS Instance with the given name
func (r *Manager) Delete(name string) (*DB, error) {
	dbInstanceInput := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(name),
		SkipFinalSnapshot:    aws.Bool(true),
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
