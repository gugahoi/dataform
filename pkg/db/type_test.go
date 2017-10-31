package db_test

import (
	"fmt"
	"testing"

	"github.com/MYOB-Technology/dataform/pkg/db"
	"github.com/aws/aws-sdk-go/service/rds"
)

func TestDB(t *testing.T) {
	testCases := []struct {
		desc, identifier, address, arn, status, instanceclass string
		masterusername, masterpassword                        string
		multiaz                                               bool
		port                                                  int64
	}{
		{
			desc:           "All Fields",
			identifier:     "some-identifier",
			address:        "some-address.com",
			port:           5432,
			arn:            "anr:1234:blah",
			multiaz:        true,
			status:         "creating",
			instanceclass:  "db.t2.micro",
			masterusername: "diavolo",
			masterpassword: "welcome",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r := rds.DBInstance{
				DBInstanceIdentifier: &tC.identifier,
				Endpoint: &rds.Endpoint{
					Address: &tC.address,
					Port:    &tC.port,
				},
				DBInstanceArn:    &tC.arn,
				MultiAZ:          &tC.multiaz,
				DBInstanceStatus: &tC.status,
				DBInstanceClass:  &tC.instanceclass,
				MasterUsername:   &tC.masterusername,
				//		StorageAllocatedGb: r.AllocatedStorage,
				//		StorageType:        r.StorageType,
				//		StorageIops:        r.Iops,
				//		StorageEncrypted:   r.StorageEncrypted,
				//		Engine:             r.Engine,
				//		EngineVersion:      r.EngineVersion,
			}

			result := db.FromDBInstance(&r)
			if result.Name != &tC.identifier {
				t.Errorf("Expected DB name %s, got %s.", tC.identifier, *result.Name)
			}
			if result.ARN != &tC.arn {
				t.Errorf("Expected DB ARN %s, got %s.", tC.arn, *result.ARN)
			}
			if result.MultiAZ != &tC.multiaz {
				t.Errorf("Expected DB MultiAZ %v, got %v.", tC.multiaz, *result.MultiAZ)
			}
			if result.Status != &tC.status {
				t.Errorf("Expected DB Status %s, got %s.", tC.status, *result.Status)
			}
			if result.DBInstanceClass != &tC.instanceclass {
				t.Errorf("Expected DB InstanceClass %s, got %s.", tC.instanceclass, *result.DBInstanceClass)
			}
			if result.Address != &tC.address {
				t.Errorf("Expected DB Address %s, got %s.", tC.address, *result.Address)
			}
			if result.Port != &tC.port {
				t.Errorf("Expected DB Port %d, got %d.", tC.port, *result.Port)
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
	az := true
	status := "status"

	db := db.DB{
		Name:    &id,
		Status:  &status,
		ARN:     &arn,
		MultiAZ: &az,
	}

	got := db.String()
	expected := fmt.Sprintf("name: %s, arn: %s", id, arn)
	if got != "name: id, arn: arn" {
		t.Errorf("Expected DB stringer to be '%v', got '%v'", expected, got)
	}
}
