package db

import (
	"fmt"
	"regexp"

	"github.com/MYOB-Technology/dataform/pkg/service"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

// RDS Instance type
type RDS struct {
	Name    string
	Status  string
	ARN     string
	AZ      string
	Address string
}

// Apply applies any allowed changed to the RDS Instance
func (r RDS) Apply() error {
	if err := validateName(r.Name); err != nil {
		return err
	}

	if r.Status == "" {
		// create
	} else {
		// update
	}
	return nil
}

// GetStatus updates the State of the DB
func (r RDS) GetStatus() (string, error) {
	// update the state of the DB Instance
	// r.sync?
	return r.Status, nil
}

func (r RDS) create(svc service.Service) error {
	svc.CreateDBInstance(r.toCreateDBInstanceInput())
	return nil
}

func (r RDS) toCreateDBInstanceInput() *rds.CreateDBInstanceInput {
	return &rds.CreateDBInstanceInput{
		DBInstanceIdentifier: aws.String(r.Name),
	}
}

// validateName checks if the name of the RDS Instance is valid
func validateName(n string) error {
	if len(n) > 63 {
		return ErrNameTooLong
	}

	// parse name (start with a-z and can include a-z0-9\-)
	re := regexp.MustCompile(`^[a-z][-a-z0-9]*$`)
	if ok := re.MatchString(n); !ok {
		return ErrInvalidName
	}
	return nil
}

// ErrInvalidName is an error that represents an invalid name
var ErrInvalidName = fmt.Errorf("db name must match [a-z][-a-z0-9]*")

// ErrNameTooLong is an error that represents a name that is more than 63 chars
var ErrNameTooLong = fmt.Errorf("db name must be between 1 and 63 chars")
