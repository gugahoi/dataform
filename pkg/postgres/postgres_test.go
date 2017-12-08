package postgres

import (
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const dc = "Don't Care"

func TestUser(t *testing.T) {
	u := &User{"foo", "baz"}
	ss := u.SQL()
	x := []string{"CREATE USER \"foo\" WITH ENCRYPTED PASSWORD 'baz'"}
	for i, s := range ss {
		if s != x[i] {
			t.Errorf("Expected %#v got %#v\n", x, s)
		}
	}
}

func TestDatabase(t *testing.T) {
	d := &Database{"foo"}
	ss := d.SQL()
	x := []string{"CREATE DATABASE \"foo\""}
	for i, s := range ss {
		if s != x[i] {
			t.Errorf("Expected %#v got %#v\n", x, s)
		}
	}
}

func TestGrantAccess(t *testing.T) {
	g := &GrantAccess{&Database{"foo"}, &User{"baz", dc}}
	ss := g.SQL()
	x := []string{
		"GRANT CONNECT ON DATABASE \"foo\" TO \"baz\"",
		"GRANT USAGE ON SCHEMA PUBLIC TO \"baz\"",
	}
	for i, s := range ss {
		if s != x[i] {
			t.Errorf("Expected %#v got %#v\n", x, s)
		}
	}
}

func TestGrantAdmin(t *testing.T) {
	g := &GrantAdmin{&Database{"foo"}, &User{"baz", dc}}
	ss := g.SQL()
	x := []string{
		"GRANT ALL PRIVILEGES ON DATABASE \"foo\" TO \"baz\" WITH GRANT OPTION",
		"GRANT ALL PRIVILEGES ON SCHEMA PUBLIC TO \"baz\" WITH GRANT OPTION",
		"GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA PUBLIC TO \"baz\" WITH GRANT OPTION",
	}
	for i, s := range ss {
		if s != x[i] {
			t.Errorf("Expected %#v got %#v\n", x, s)
		}
	}
}

func TestGrantRead(t *testing.T) {
	g := &GrantRead{&User{"foo", dc}}
	ss := g.SQL()
	x := []string{
		"GRANT SELECT ON ALL TABLES IN SCHEMA PUBLIC TO \"foo\"",
		"ALTER DEFAULT PRIVILEGES IN SCHEMA PUBLIC GRANT SELECT ON TABLES TO \"foo\"",
	}
	for i, s := range ss {
		if s != x[i] {
			t.Errorf("Expected %#v got %#v\n", x, s)
		}
	}
}

func TestGrantWrite(t *testing.T) {
	g := &GrantWrite{&User{"foo", dc}}
	ss := g.SQL()
	x := []string{
		"GRANT SELECT,INSERT,UPDATE,DELETE,REFERENCES ON ALL TABLES IN SCHEMA PUBLIC TO \"foo\"",
		"GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA PUBLIC TO \"foo\"",
		"ALTER DEFAULT PRIVILEGES IN SCHEMA PUBLIC GRANT SELECT,INSERT,UPDATE,DELETE,REFERENCES ON TABLES TO \"foo\"",
	}
	for i, s := range ss {
		if s != x[i] {
			t.Errorf("Expected %#v got %#v\n", x, s)
		}
	}
}

func TestExecCreateDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to mock: %v\n", err)
	}

	ex := "^CREATE DATABASE \"db\""
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "REVOKE ALL PRIVILEGES ON SCHEMA PUBLIC FROM PUBLIC CASCADE"
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "REVOKE ALL PRIVILEGES ON DATABASE \"db\" FROM PUBLIC CASCADE"
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "CREATE USER \"db-admin\" WITH ENCRYPTED PASSWORD 'admin'"
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "GRANT ALL PRIVILEGES ON DATABASE \"db\" TO \"db-admin\" WITH GRANT OPTION"
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "GRANT ALL PRIVILEGES ON SCHEMA PUBLIC TO \"db-admin\" WITH GRANT OPTION"
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA PUBLIC TO \"db-admin\" WITH GRANT OPTION"
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "CREATE USER \"db-writer\" WITH ENCRYPTED PASSWORD 'writer'"
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "CREATE USER \"db-reader\" WITH ENCRYPTED PASSWORD 'reader'"
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	dd := &DatabaseDescriptor{
		Host:     "test",
		Port:     0,
		Database: &Database{"db"},
		Admin:    &User{"db-admin", "admin"},
		Writer:   &User{"db-writer", "writer"},
		Reader:   &User{"db-reader", "reader"},
	}

	c := &Conn{"dummy", 0, &User{dc, dc}, db}

	if err := c.execCreateDatabase(dd); err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestExecCreateDatabaseUserPrivs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to mock: %v\n", err)
	}

	ex := "REVOKE ALL PRIVILEGES ON SCHEMA PUBLIC FROM PUBLIC CASCADE"
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "REVOKE ALL PRIVILEGES ON DATABASE \"db\" FROM PUBLIC CASCADE"
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "GRANT CONNECT ON DATABASE \"db\" TO \"db-writer\""
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "GRANT USAGE ON SCHEMA PUBLIC TO \"db-writer\""
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "GRANT SELECT,INSERT,UPDATE,DELETE,REFERENCES ON ALL TABLES IN SCHEMA PUBLIC TO \"db-writer\""
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA PUBLIC TO \"db-writer\""
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "ALTER DEFAULT PRIVILEGES IN SCHEMA PUBLIC GRANT SELECT,INSERT,UPDATE,DELETE,REFERENCES ON TABLES TO \"db-writer\""
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "GRANT CONNECT ON DATABASE \"db\" TO \"db-reader\""
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "GRANT USAGE ON SCHEMA PUBLIC TO \"db-reader\""
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "GRANT SELECT ON ALL TABLES IN SCHEMA PUBLIC TO \"db-reader\""
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	ex = "ALTER DEFAULT PRIVILEGES IN SCHEMA PUBLIC GRANT SELECT ON TABLES TO \"db-reader\""
	mock.ExpectExec(ex).WillReturnResult(sqlmock.NewResult(0, 0))

	dd := &DatabaseDescriptor{
		Host:     "test",
		Port:     0,
		Database: &Database{"db"},
		Admin:    &User{"db-admin", "admin"},
		Writer:   &User{"db-writer", "writer"},
		Reader:   &User{"db-reader", "reader"},
	}

	c := &Conn{"dummy", 0, &User{dc, dc}, db}

	if err := c.execCreateDatabaseUserPrivs(dd); err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
