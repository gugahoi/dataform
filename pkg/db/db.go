package db

// DB is an interface for databases
type DB interface {
	Create() error
	Delete() error
	Update() error
}

var (
	// StatusOK is an RDS create status
	StatusOK = "created"
	// StatusCreating is an RDS creating status
	StatusCreating = "creating"
	// StatusDeleting is an RDS deleting status
	StatusDeleting = "deleting"
	// StatusBackingUp is an RDS backup status
	StatusBackingUp = "backing-up"
	// StatusAvailable is an RDS available status
	StatusAvailable = "available"
)
