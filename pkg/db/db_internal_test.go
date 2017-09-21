package db

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

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
	cases := []struct {
		name string
		err  error
	}{
		{name: "gohan", err: nil},
		{name: "goku", err: fmt.Errorf("goku error")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := mockRdsSvc{
				err: tt.err,
			}
			rds := NewManager(svc)

			db, err := rds.Create(tt.name)
			if err != tt.err {
				t.Errorf("Expected error to be %v, got %v", tt.err, err)
			}

			if db != nil {
				t.Errorf("Expected db to be %v, got %v", nil, db)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{name: "gohan", err: nil},
		{name: "goku", err: fmt.Errorf("Goku Error")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := mockRdsSvc{
				err: tt.err,
			}
			rds := NewManager(svc)

			db, err := rds.Delete(tt.name)
			if err != tt.err {
				t.Errorf("Expected error to be %v, got %v", tt.err, err)
			}

			if db != nil {
				t.Errorf("Expected db to be %v, got %v", nil, db)
			}
		})
	}
}

func TestStatus(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{name: "gohan", err: nil},
		{name: "goku", err: fmt.Errorf("Goku Error")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := mockRdsSvc{
				err: tt.err,
			}
			rds := NewManager(svc)

			db, err := rds.Stat(tt.name)
			if err != tt.err {
				t.Errorf("Expected error to be %v, got %v", tt.err, err)
			}

			if db != nil {
				t.Errorf("Expected db to be %v, got %v", nil, db)
			}
		})
	}
}

func TestStatAll(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{name: "gohan", err: nil},
		{name: "goku", err: fmt.Errorf("Goku Error")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := mockRdsSvc{
				err: tt.err,
			}
			rds := NewManager(svc)

			db, err := rds.StatAll()
			if err != tt.err {
				t.Errorf("Expected error to be %v, got %v", tt.err, err)
			}

			if db != nil {
				t.Errorf("Expected db to be %v, got %v", nil, db)
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
