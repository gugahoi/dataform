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
func (r *Manager) Create(name string) (*rds.CreateDBInstanceOutput, error) {
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

	return result, nil
}

// Delete an RDS Instance with the given name
func (r *Manager) Delete(name string) (*rds.DeleteDBInstanceOutput, error) {
	dbInstanceInput := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(name),
		SkipFinalSnapshot:    aws.Bool(true),
	}

	result, err := r.Client.DeleteDBInstance(dbInstanceInput)
	if err != nil {
		return nil, err
	}

	return result, nil
}
