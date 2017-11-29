package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
)

// Tag encapsulates an aws resource tag key/value
type Tag struct {
	Key   *string
	Value *string
}

// Profiles enum
const (
	Production = iota
	Development
)

// DB InstanceParam type
type InstanceParams struct {
	Name                       *string
	Status                     *string
	Engine                     *string
	EngineVersion              *string
	ARN                        *string
	CopyTagsToSnapshot         *bool
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

// ProfileInstanceParams these can change based on the profile
type ProfileInstanceParams struct {
	MultiAZ               *bool
	BackupRetentionPeriod *int64
}

// DB Instance Type
type DB struct {
	InstanceParams
	ProfileInstanceParams
}

// FromDBInstance converts an *rds.DBInstance type to *DB type
func FromDBInstance(r *rds.DBInstance) *DB {
	var Params = InstanceParams{
		ARN:                r.DBInstanceArn,
		Name:               r.DBInstanceIdentifier,
		Status:             r.DBInstanceStatus,
		CopyTagsToSnapshot: r.CopyTagsToSnapshot,
		DBInstanceClass:    r.DBInstanceClass,
		MasterUsername:     r.MasterUsername,
		StorageAllocatedGB: r.AllocatedStorage,
		StorageType:        r.StorageType,
		StorageIops:        r.Iops,
		StorageEncrypted:   r.StorageEncrypted,
		Engine:             r.Engine,
		EngineVersion:      r.EngineVersion,
	}

	var ProfileParams = ProfileInstanceParams{
		MultiAZ:               r.MultiAZ,
		BackupRetentionPeriod: r.BackupRetentionPeriod,
	}

	var db = &DB{
		InstanceParams:        Params,
		ProfileInstanceParams: ProfileParams,
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
	if r.Endpoint != nil && r.Endpoint.Port != nil {
		db.Port = r.Endpoint.Port
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
