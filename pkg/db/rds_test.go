package db_test

import (
	"testing"

	"github.com/MYOB-Technology/dataform/pkg/db"
)

func TestApply(t *testing.T) {
	testCases := []struct {
		desc   string
		name   string
		status string
		err    error
	}{
		{
			desc:   "new rds with valid name",
			name:   "valid-name",
			status: "",
			err:    nil,
		},
		{
			desc:   "new rds with invalid name",
			name:   "123-asd",
			status: "",
			err:    db.ErrInvalidName,
		},
		{
			desc:   "existing rds with valid name",
			name:   "some-rds",
			status: db.StatusAvailable,
			err:    nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			d := db.RDS{
				Name:   tC.name,
				Status: tC.status,
			}

			if err := d.Apply(); err != tC.err {
				t.Errorf("expected err to be '%v', got '%v'", tC.err, err)
			}
		})
	}
}

func TestGetStatus(t *testing.T) {
	testCases := []struct {
		desc   string
		status string
		err    error
	}{
		{
			desc:   "simple status update",
			status: db.StatusDeleting,
			err:    nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			d := db.RDS{
				Name:   "valid-name",
				Status: tC.status,
			}

			status, err := d.GetStatus()
			if err != tC.err {
				t.Errorf("expected err to be '%v', got '%v'", tC.err, err)
			}

			if status != tC.status {
				t.Errorf("expected status to be '%v', got '%v'", tC.status, status)
			}
		})
	}
}
