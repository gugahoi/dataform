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
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			DBInstance := rds.DBInstance{
				AvailabilityZone:     &tt.az,
				DBInstanceIdentifier: &tt.name,
				DBInstanceArn:        &tt.arn,
				Endpoint:             &rds.Endpoint{},
			}
			expectedDB := FromDBInstance(&DBInstance)

			svc := mockRdsSvc{
				err: tt.err,
				CreateDBInstanceOutput: &rds.CreateDBInstanceOutput{
					DBInstance: &DBInstance,
				},
			}
			rds := NewManager(svc)

			db, err := rds.Create(tt.name)
			if err != tt.err {
				t.Errorf("Expected error to be %v, got %v", tt.err, err)
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
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			DBInstance := rds.DBInstance{
				AvailabilityZone:     &tt.az,
				DBInstanceIdentifier: &tt.name,
				DBInstanceArn:        &tt.arn,
			}
			expectedDB := FromDBInstance(&DBInstance)

			svc := mockRdsSvc{
				err: tt.err,
				DeleteDBInstanceOutput: &rds.DeleteDBInstanceOutput{
					DBInstance: &DBInstance,
				},
			}
			rds := NewManager(svc)

			db, err := rds.Delete(tt.name)
			if err != tt.err {
				t.Errorf("Expected error to be %v, got %v", tt.err, err)
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
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			DBInstances := []*rds.DBInstance{}
			for i := 0; i < tt.instances; i++ {
				DBInstances = append(DBInstances, &rds.DBInstance{})
			}

			expectedDBs := FromDBInstances(DBInstances)

			svc := mockRdsSvc{
				err: tt.err,
				DescribeDBInstancesOutput: &rds.DescribeDBInstancesOutput{
					DBInstances: DBInstances,
				},
			}
			rds := NewManager(svc)

			result, err := rds.List()
			if err != tt.err {
				t.Errorf("Expected error to be %v, got %v", tt.err, err)
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
	DeleteDBInstanceOutput    *rds.DeleteDBInstanceOutput
	DescribeDBInstancesOutput *rds.DescribeDBInstancesOutput
	err                       error
}

func (m mockRdsSvc) CreateDBInstance(input *rds.CreateDBInstanceInput) (*rds.CreateDBInstanceOutput, error) {
	return m.CreateDBInstanceOutput, m.err
}

func (m mockRdsSvc) DeleteDBInstance(input *rds.DeleteDBInstanceInput) (*rds.DeleteDBInstanceOutput, error) {
	return m.DeleteDBInstanceOutput, m.err
}

func (m mockRdsSvc) DescribeDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
	return m.DescribeDBInstancesOutput, m.err
}
