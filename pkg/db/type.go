package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
)

// DB Instance type
type DB struct {
	Name    *string
	Status  *string
	ARN     *string
	AZ      *string
	Address *string
}

// FromDBInstance converts an *rds.DBInstance type to *DB type
func FromDBInstance(r *rds.DBInstance) *DB {
	// var address *string
	// if r.Endpoint != nil && r.Endpoint.Address != nil {
	// 	address = r.Endpoint.Address
	// }

	return &DB{
		ARN:    r.DBInstanceArn,
		Name:   r.DBInstanceIdentifier,
		Status: r.DBInstanceStatus,
		// AZ:     r.AvailabilityZone,
		// Address: address,
	}
}

// FromDBInstances converts a slice of *rds.DBInstance to a slice of *DB
func FromDBInstances(r []*rds.DBInstance) []*DB {
	var DBs []*DB
	for _, instance := range r {
		DBs = append(DBs, FromDBInstance(instance))
	}

	return DBs
}

// String representation of DB
func (d *DB) String() string {
	return fmt.Sprintf("name: %s, arn: %s", *d.Name, *d.ARN)
}
