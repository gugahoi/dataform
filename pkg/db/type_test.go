package db_test

import (
	"fmt"
	"testing"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
)

func TestDB(t *testing.T) {
	testCases := []struct {
		desc, identifier, address, arn, az, status string
	}{
		{
			desc:       "All Fields",
			identifier: "some-identifier",
			address:    "some-address.com",
			arn:        "anr:1234:blah",
			az:         "southeast-2b,",
			status:     "creating",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r := rds.DBInstance{
				DBInstanceIdentifier: &tC.identifier,
				Endpoint: &rds.Endpoint{
					Address: &tC.address,
				},
				DBInstanceArn:    &tC.arn,
				AvailabilityZone: &tC.az,
				DBInstanceStatus: &tC.status,
			}

			result := db.FromDBInstance(&r)
			if result.Name != &tC.identifier {
				t.Errorf("Expected DB name %s, got %s.", tC.identifier, *result.Name)
			}
			if result.ARN != &tC.arn {
				t.Errorf("Expected DB ARN %s, got %s.", tC.arn, *result.ARN)
			}
			if result.AZ != &tC.az {
				t.Errorf("Expected DB AZ %s, got %s.", tC.az, *result.AZ)
			}
			if result.Status != &tC.status {
				t.Errorf("Expected DB Status %s, got %s.", tC.status, *result.Status)
			}
			if result.Address != &tC.address {
				t.Errorf("Expected DB Address %s, got %s.", tC.address, *result.Address)
			}
		})
	}
}

func TestFromDBInstances(t *testing.T) {
	testCases := []struct {
		desc string
		size int
	}{
		{
			desc: "No Instances",
			size: 0,
		},
		{
			desc: "Ten Instances",
			size: 10,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			DBs := []*rds.DBInstance{}
			for i := 0; i < tC.size; i++ {
				arn := fmt.Sprintf("arn-%d", i)
				id := fmt.Sprintf("id-%d", i)
				az := fmt.Sprintf("az-%d", i)
				status := fmt.Sprintf("status-%d", i)
				DBs = append(DBs, &rds.DBInstance{
					DBInstanceArn:        &arn,
					DBInstanceIdentifier: &id,
					AvailabilityZone:     &az,
					DBInstanceStatus:     &status,
				})
			}
			result := db.FromDBInstances(DBs)
			if len(result) != tC.size {
				t.Errorf("Expected %d instances, got %d", tC.size, len(result))
			}
		})
	}
}

func TestString(t *testing.T) {
	arn := "arn"
	id := "id"
	az := "az"
	status := "status"

	db := db.DB{
		Name:   &id,
		Status: &status,
		ARN:    &arn,
		AZ:     &az,
	}

	got := db.String()
	expected := fmt.Sprintf("name: %s, arn: %s", id, arn)
	if got != "name: id, arn: arn" {
		t.Errorf("Expected DB stringer to be '%v', got '%v'", expected, got)
	}
}
