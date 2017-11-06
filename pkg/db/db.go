package db

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	time "time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

var (
	// Defaults are rds instance defaults
	Defaults = &DB{
		DBInstanceClass:    aws.String("db.t2.small"),
		CopyTagsToSnapshot: aws.Bool(true),
		Engine:             aws.String("postgres"),
		EngineVersion:      aws.String("9.6.3"),
		MultiAZ:            aws.Bool(false),
		Port:               aws.Int64(5432),
		StorageAllocatedGB: aws.Int64(5),
		StorageEncrypted:   aws.Bool(true),
		StorageType:        aws.String("gp2"),
	}

	errInvalidUsernamePassword           = fmt.Errorf("username and password cannot be empty")
	errDbNameMissing                     = fmt.Errorf("error: required DB field Name is missing")
	errDbMasterUsernameMissing           = fmt.Errorf("error: required DB field MasterUsername is missing")
	errDbMasterUserPasswordMissing       = fmt.Errorf("error: required DB field MasterUserPassword is missing")
	errStateTransitionedToErrorCondition = fmt.Errorf("error: db transitioned to error condition")
)

// Manager uses a svc to talk to AWS RDS
type Manager struct {
	Client  rdsiface.RDSAPI
	wait    sync.WaitGroup
	signals chan os.Signal
	stop    chan struct{}
}

// NewManager returns a pointer to a Manager struct.
// The supplied svc is used to make calls to AWS RDS Service.
func NewManager(svc rdsiface.RDSAPI) *Manager {
	stop := make(chan struct{})
	signals := make(chan os.Signal, 1)
	return &Manager{
		Client:  svc,
		wait:    sync.WaitGroup{},
		stop:    stop,
		signals: signals,
	}
}

// SigHandler handles signals. Run as a goroutine
func (r *Manager) SigHandler() {
	r.wait.Add(1)
	defer r.wait.Done()
	signal.Notify(r.signals,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	sig := <-r.signals
	fmt.Printf("signal received: %v\n", sig)
	r.Stop()
	return
}

// SetShutdownChannel allows the use of an external channel to signal closure
func (r *Manager) SetShutdownChannel(stop chan struct{}) {
	r.stop = stop
}

// Stop will tell any goroutines  to shutdown
func (r *Manager) Stop() {
	close(r.stop)
	r.wait.Wait()
	// Release all remaining resources.
	r.stop = nil
}

// Create an RDS Instance from a supplied DB object
func (r *Manager) Create(db *DB) (*DB, error) {
	database, err := validateDBInput(db)
	if err != nil {
		return nil, err
	}
	dbInput := &rds.CreateDBInstanceInput{
		AllocatedStorage:     database.StorageAllocatedGB,
		DBInstanceClass:      database.DBInstanceClass,
		DBInstanceIdentifier: database.Name,
		Engine:               database.Engine,
		EngineVersion:        database.EngineVersion,
		MasterUserPassword:   database.MasterUserPassword,
		MasterUsername:       database.MasterUsername,
		MultiAZ:              database.MultiAZ,
		Port:                 database.Port,
		DBSubnetGroupName:    database.SubnetGroupName,
		VpcSecurityGroupIds:  database.SecurityGroups,
		StorageEncrypted:     database.StorageEncrypted,
		StorageType:          database.StorageType,
		Iops:                 database.StorageIops,
	}

	if database.KMSKeyArn != nil {
		dbInput.KmsKeyId = database.KMSKeyArn
		database.StorageEncrypted = aws.Bool(true)
	}
	if database.PreferredBackupWindow != nil {
		dbInput.PreferredBackupWindow = database.PreferredBackupWindow
	}
	if database.PreferredMaintenanceWindow != nil {
		dbInput.PreferredMaintenanceWindow = database.PreferredMaintenanceWindow
	}
	if database.Tags != nil {
		tags := make([]*rds.Tag, 0, 10)
		for _, v := range database.Tags {
			tags = append(tags, &rds.Tag{
				Key:   v.Key,
				Value: v.Value,
			})
		}
		dbInput.Tags = tags
	}

	result, err := r.Client.CreateDBInstance(dbInput)
	if err != nil {
		return nil, err
	}

	return FromDBInstance(result.DBInstance), nil
}

func validateDBInput(db *DB) (*DB, error) {
	if db.Name == nil {
		return nil, errDbNameMissing
	}
	if db.MasterUsername == nil {
		return nil, errDbMasterUsernameMissing
	}
	if len(*db.MasterUsername) == 0 {
		return nil, errDbMasterUsernameMissing
	}
	if db.MasterUserPassword == nil {
		return nil, errDbMasterUserPasswordMissing
	}
	if len(*db.MasterUserPassword) == 0 {
		return nil, errDbMasterUserPasswordMissing
	}

	// cannot override this one for now
	db.Engine = Defaults.Engine
	db.CopyTagsToSnapshot = Defaults.CopyTagsToSnapshot

	if db.MultiAZ == nil {
		db.MultiAZ = Defaults.MultiAZ
	}
	if db.DBInstanceClass == nil {
		db.DBInstanceClass = Defaults.DBInstanceClass
	}
	if db.EngineVersion == nil {
		db.EngineVersion = Defaults.EngineVersion
	}
	if db.StorageAllocatedGB == nil {
		db.StorageAllocatedGB = Defaults.StorageAllocatedGB
	}
	if db.StorageType == nil {
		db.StorageType = Defaults.StorageType
	}
	if db.StorageEncrypted == nil {
		db.StorageEncrypted = Defaults.StorageEncrypted
	}

	return db, nil
}

// Delete an RDS Instance with the given name
func (r *Manager) Delete(name string) (*DB, error) {
	now := time.Now()
	snapshotID := fmt.Sprintf("%s-%s", name, now.Format("20060102150405"))
	dbInstanceInput := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier:      aws.String(name),
		FinalDBSnapshotIdentifier: aws.String(snapshotID),
		SkipFinalSnapshot:         aws.Bool(false),
	}

	result, err := r.Client.DeleteDBInstance(dbInstanceInput)
	if err != nil {
		return nil, err
	}

	return FromDBInstance(result.DBInstance), nil
}

// Stat returns the status of an RDS Instance
func (r *Manager) Stat(name string) (*DB, error) {
	dbInstanceInput := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(name),
	}

	result, err := r.Client.DescribeDBInstances(dbInstanceInput)
	if err != nil {
		return nil, err
	}

	if len(result.DBInstances) == 0 {
		return nil, nil
	}

	return FromDBInstance(result.DBInstances[0]), nil
}

// List returns the status of all RDS Instances
func (r *Manager) List() ([]*DB, error) {
	dbInstanceInput := &rds.DescribeDBInstancesInput{}

	result, err := r.Client.DescribeDBInstances(dbInstanceInput)
	if err != nil {
		return nil, err
	}

	return FromDBInstances(result.DBInstances), nil
}

// State is used to return whether DB state is finalised or not
type State struct {
	Final  bool
	Status string
	Err    error
}

// WaitForFinalState will block until the requested instance is in a known final state
func (r *Manager) WaitForFinalState(dbname string, pollInterval time.Duration, pollTimeout time.Duration) <-chan State {
	result := make(chan State)
	go func() {
		timeout := time.After(pollTimeout * time.Second)
		tick := time.Tick(pollInterval * time.Second)
		r.wait.Add(1)
		defer close(result)
		defer r.wait.Done()
		for {
			select {
			case <-r.stop:
				result <- State{
					Final:  false,
					Status: "",
					Err:    nil,
				}
				return
			case <-timeout:
				result <- State{
					Final:  false,
					Status: "timeout",
					Err:    fmt.Errorf("error timed out polling for db final state: %s", dbname),
				}
				return
			case <-tick:
				db, err := r.Stat(dbname)
				if err != nil {
					result <- State{
						Final:  true,
						Status: StatusDeleted,
						Err:    err,
					}
					return
				}
				status := r.IsFinalState(db)
				result <- status
				if status.Final {
					return
				}
			}
		}
	}()
	return result
}

// IsFinalState checks whether the current rds state is final, transitioning, or in an error state
func (r *Manager) IsFinalState(db *DB) State {
	if db == nil {
		return State{
			Final:  true,
			Status: StatusDeleted,
			Err:    nil,
		}
	}
	if FinalStates[*db.Status] {
		return State{
			Final:  true,
			Status: *db.Status,
			Err:    nil,
		}
	}
	if TransitioningStates[*db.Status] {
		return State{
			Final:  false,
			Status: *db.Status,
			Err:    nil,
		}
	}
	// if we get here, all other states are error conditions
	return State{
		Final: true,
		Err:   errStateTransitionedToErrorCondition,
	}
}

// generateRandomString receives a size and a string of allowed characters and generates a random string of the given size
func generateRandomString(strlen int, allowedChars string, t clock) string {
	rsource := rand.New(rand.NewSource(t.Now().UnixNano()))
	result := make([]byte, strlen)
	for i := range result {
		result[i] = allowedChars[rsource.Intn(len(allowedChars))]
	}
	return string(result)
}

// GenerateRandomPassword receives a size and generates a random password of that size
func (r *Manager) GenerateRandomPassword(strlen int) string {
	const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFHIJKLMNOPQRSTUVWXYZ0123456789!#$%^&*()-+="
	return generateRandomString(strlen, allowedChars, actualClock{})
}

// GenerateRandomUsername receives a size and generates a random username of that size
func (r *Manager) GenerateRandomUsername(strlen int) string {
	const allowedChars = "abcdefghijklmnopqrstuvwxyz"
	return generateRandomString(strlen, allowedChars, actualClock{})
}

// clock allows us to mock out time.Now in our tests
type clock interface {
	Now() time.Time
}

type actualClock struct{}

func (actualClock) Now() time.Time {
	return time.Now()
}
