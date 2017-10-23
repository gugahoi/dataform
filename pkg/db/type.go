package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
)

type Tag struct {
	Key   *string
	Value *string
}

// DB Instance type
type DB struct {
	Name                       *string
	Status                     *string
	Engine                     *string
	EngineVersion              *string
	ARN                        *string
	CopyTagsToSnapshot         *bool
	MultiAZ                    *bool
	Address                    *string
	DBInstanceClass            *string
	KMSKeyArn                  *string
	Port                       *int64
	SubnetGroupName            *string
	SecurityGroups             []*string
	MasterUserPassword         *string
	MasterUsername             *string
	StorageAllocatedGB         *int64
	StorageType                *string
	StorageIops                *int64
	StorageEncrypted           *bool
	PreferredBackupWindow      *string
	PreferredMaintenanceWindow *string
	Tags                       []*Tag
}

// FromDBInstance converts an *rds.DBInstance type to *DB type
func FromDBInstance(r *rds.DBInstance) *DB {
	db := &DB{
		ARN:                r.DBInstanceArn,
		Name:               r.DBInstanceIdentifier,
		Status:             r.DBInstanceStatus,
		CopyTagsToSnapshot: r.CopyTagsToSnapshot,
		MultiAZ:            r.MultiAZ,
		DBInstanceClass:    r.DBInstanceClass,
		// SecurityGroups:     r.VpcSecurityGroups,
		MasterUsername:     r.MasterUsername,
		StorageAllocatedGB: r.AllocatedStorage,
		StorageType:        r.StorageType,
		StorageIops:        r.Iops,
		StorageEncrypted:   r.StorageEncrypted,
		Engine:             r.Engine,
		EngineVersion:      r.EngineVersion,
	}
	if r.KmsKeyId != nil {
		db.KMSKeyArn = r.KmsKeyId
	}
	if r.DBSubnetGroup != nil && r.DBSubnetGroup.DBSubnetGroupName != nil {
		db.SubnetGroupName = r.DBSubnetGroup.DBSubnetGroupName
	}
	if r.Endpoint != nil && r.Endpoint.Address != nil {
		db.Address = r.Endpoint.Address
	}
	return db
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
