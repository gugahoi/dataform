package db

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

var cases = []struct {
	name, az, arn string
	instances     int
	err           error
}{
	{name: "Happy Path", az: "gohan-az", arn: "here-is-the-arn", instances: 2, err: nil},
	{name: "Sad Path", err: fmt.Errorf("Goku Error")},
}

func TestNewRds(t *testing.T) {
	svc := rds.New(session.New(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	}))

	rds := NewManager(svc)

	if rds == nil {
		t.Errorf("Name does not match")
	}
}

func TestCreate(t *testing.T) {
	var cases = []struct {
		name, az, arn      string
		username, password string
		instances          int
		err                error
	}{
		{name: "Happy Path", username: "trunks", password: "bulma", az: "gohan-az", arn: "here-is-the-arn", instances: 2, err: nil},
		{name: "Empty Username", username: "", password: "bulma", err: errInvalidUsernamePassword},
		{name: "Empty Password", username: "trunks", password: "", err: errInvalidUsernamePassword},
	}

	for _, tC := range cases {
		t.Run(tC.name, func(t *testing.T) {
			DBInstance := rds.DBInstance{
				AvailabilityZone:     &tC.az,
				DBInstanceIdentifier: &tC.name,
				DBInstanceArn:        &tC.arn,
				Endpoint:             &rds.Endpoint{},
			}
			expectedDB := FromDBInstance(&DBInstance)

			svc := mockRdsSvc{
				err: tC.err,
				CreateDBInstanceOutput: &rds.CreateDBInstanceOutput{
					DBInstance: &DBInstance,
				},
			}
			rds := NewManager(svc)

			db, err := rds.Create(tC.name, tC.username, tC.password)
			if err != tC.err {
				t.Errorf("Expected error to be %v, got %v", tC.err, err)
			}

			if db != nil {
				if db.ARN != expectedDB.ARN {
					t.Errorf("Expected db arn to be %v, got %v", expectedDB, db)
				}
				if db.Name != expectedDB.Name {
					t.Errorf("Expected db name to be %v, got %v", expectedDB, db)
				}
				if db.AZ != expectedDB.AZ {
					t.Errorf("Expected db AZ to be %v, got %v", expectedDB, db)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	for _, tC := range cases {
		t.Run(tC.name, func(t *testing.T) {
			DBInstance := rds.DBInstance{
				AvailabilityZone:     &tC.az,
				DBInstanceIdentifier: &tC.name,
				DBInstanceArn:        &tC.arn,
			}
			expectedDB := FromDBInstance(&DBInstance)

			svc := mockRdsSvc{
				err: tC.err,
				DeleteDBInstanceOutput: &rds.DeleteDBInstanceOutput{
					DBInstance: &DBInstance,
				},
			}
			rds := NewManager(svc)

			db, err := rds.Delete(tC.name)
			if err != tC.err {
				t.Errorf("Expected error to be %v, got %v", tC.err, err)
			}

			if db != nil {
				if db.ARN != expectedDB.ARN {
					t.Errorf("Expected db ARN to be %v, got %v", expectedDB.ARN, db.ARN)
				}
				if db.Name != expectedDB.Name {
					t.Errorf("Expected db name to be %v, got %v", expectedDB.Name, db.Name)
				}
				if db.AZ != expectedDB.AZ {
					t.Errorf("Expected db AZ to be %v, got %v", expectedDB.AZ, db.AZ)
				}
			}
		})
	}
}

func TestStatus(t *testing.T) {
	name := "db-stating"
	arn := "arn:123:123:rds:db-stating"
	var expectedErr error
	DBInstance := rds.DBInstance{
		DBInstanceIdentifier: &name,
		DBInstanceArn:        &arn,
	}
	expectedDB := FromDBInstance(&DBInstance)

	svc := mockRdsSvc{
		err: expectedErr,
		DescribeDBInstancesOutput: &rds.DescribeDBInstancesOutput{
			DBInstances: []*rds.DBInstance{
				&DBInstance,
			},
		},
	}

	rds := NewManager(svc)

	db, err := rds.Stat(name)
	if err != expectedErr {
		t.Errorf("Expected error to be %v, got %v", expectedErr, err)
	}

	if db.Name != expectedDB.Name {
		t.Errorf("Expected db name to be %v, got %v", expectedDB.Name, db.Name)
	}

	if db.ARN != expectedDB.ARN {
		t.Errorf("Expected db ARN to be %v, got %v", expectedDB.ARN, db.ARN)
	}
}

func TestList(t *testing.T) {
	for _, tC := range cases {
		t.Run(tC.name, func(t *testing.T) {
			DBInstances := []*rds.DBInstance{}
			for i := 0; i < tC.instances; i++ {
				DBInstances = append(DBInstances, &rds.DBInstance{})
			}

			expectedDBs := FromDBInstances(DBInstances)

			svc := mockRdsSvc{
				err: tC.err,
				DescribeDBInstancesOutput: &rds.DescribeDBInstancesOutput{
					DBInstances: DBInstances,
				},
			}
			rds := NewManager(svc)

			result, err := rds.List()
			if err != tC.err {
				t.Errorf("Expected error to be %v, got %v", tC.err, err)
			}

			if len(result) != len(expectedDBs) {
				t.Errorf("Expected %d results, got %d", len(expectedDBs), len(result))
			}
		})
	}
}

type mockRdsSvc struct {
	rdsiface.RDSAPI
	CreateDBInstanceOutput    *rds.CreateDBInstanceOutput
	CreateMasterUsername      *string
	CreateMasterPassword      *string
	DeleteDBInstanceOutput    *rds.DeleteDBInstanceOutput
	DescribeDBInstancesOutput *rds.DescribeDBInstancesOutput
	err                       error
}

func (m mockRdsSvc) CreateDBInstance(input *rds.CreateDBInstanceInput) (*rds.CreateDBInstanceOutput, error) {
	m.CreateMasterUsername = input.MasterUsername
	m.CreateMasterPassword = input.MasterUserPassword
	return m.CreateDBInstanceOutput, m.err
}

func (m mockRdsSvc) DeleteDBInstance(input *rds.DeleteDBInstanceInput) (*rds.DeleteDBInstanceOutput, error) {
	return m.DeleteDBInstanceOutput, m.err
}

func (m mockRdsSvc) DescribeDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
	return m.DescribeDBInstancesOutput, m.err
}
