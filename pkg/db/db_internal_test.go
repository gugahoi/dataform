package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

var cases = []struct {
	name, arn string
	multiaz   bool
	instances int
	err       error
}{
	{name: "Happy Path", multiaz: true, arn: "here-is-the-arn", err: nil},
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
	defaultInstanceType := "db.t2.micro"
	var cases = []struct {
		name, arn, masterusername, masterpassword string
		multiaz                                   bool
		username, password, dbinstanceclass       string
		instances                                 int
		err                                       error
	}{
		{name: "Happy Path", username: "trunks", password: "bulma", multiaz: true, dbinstanceclass: defaultInstanceType, arn: "here-is-the-arn", err: nil},
		{name: "Empty Username", username: "", password: "bulma", multiaz: true, err: errDbMasterUsernameMissing},
		{name: "Empty Password", username: "trunks", password: "", multiaz: true, err: errDbMasterUserPasswordMissing},
		{name: "Missing Username", password: "bulma", multiaz: true, err: errDbMasterUsernameMissing},
		{name: "Missing Password", username: "trunks", multiaz: true, err: errDbMasterUserPasswordMissing},
		{name: "Invalid Instance Type", username: "trunks", multiaz: true, dbinstanceclass: "db.t2.duff", err: errDbMasterUserPasswordMissing},
	}

	for _, tC := range cases {
		t.Run(tC.name, func(t *testing.T) {
			DBInput := &DB{
				Name:               &tC.name,
				MasterUsername:     &tC.username,
				MasterUserPassword: &tC.password,
				MultiAZ:            &tC.multiaz,
				DBInstanceClass:    &tC.dbinstanceclass,
			}
			DBInstance := rds.DBInstance{
				MultiAZ:              &tC.multiaz,
				DBInstanceIdentifier: &tC.name,
				DBInstanceArn:        &tC.arn,
				DBInstanceClass:      &defaultInstanceType,
				MasterUsername:       &tC.username,
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

			db, err := rds.Create(DBInput)
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
				if db.MultiAZ != expectedDB.MultiAZ {
					t.Errorf("Expected db MultiAZ to be %v, got %v", expectedDB, db)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	for _, tC := range cases {
		t.Run(tC.name, func(t *testing.T) {
			DBInstance := rds.DBInstance{
				MultiAZ:              &tC.multiaz,
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
				if db.MultiAZ != expectedDB.MultiAZ {
					t.Errorf("Expected db MultiAZ to be %v, got %v", expectedDB.MultiAZ, db.MultiAZ)
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

func TestGenerateRandomString(t *testing.T) {
	testCases := []struct {
		desc     string
		size     int
		allowed  string
		expected string
	}{
		{
			desc: "0 size", size: 0, allowed: "abc", expected: "",
		},
		{
			desc: "Valid Size and Some Chars", size: 10, allowed: "abc", expected: "accbacccac",
		},
		{
			desc: "Valid size with symbols", size: 8, allowed: "!@#$%", expected: "@@@%@#@@",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := generateRandomString(tC.size, tC.allowed, mockClock{})
			if got != tC.expected {
				t.Errorf("expected string to be %v, got %v", tC.expected, got)
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

// mocked clock
type mockClock struct{}

func (mockClock) Now() time.Time {
	return time.Unix(123456, 0)
}
