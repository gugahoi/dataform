// Package postgres creates databases and users on a PostgreSQL server.
//
// Example:
// conn, err := NewConn(5432, "host", "user", "password")
// err := conn.CreateDatabase("foobaz")
//
package postgres

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	_ "github.com/lib/pq"
)

// DatabaseDescriptor describes a created database.
type DatabaseDescriptor struct {
	Host     string
	Port     int
	Database *Database
	Admin    *User
	Writer   *User
	Reader   *User
}

// Conn is a connection to a Postgres server.
type Conn struct {
	Host string
	Port int
	User *User
	DB   *sql.DB
}

// NewConn creates a new Conn with the supplied host and user details.
func NewConn(port int, host, user, password string) (*Conn, error) {
	dsn := fmt.Sprintf(
		"user=%v password='%v' host=%v port=%v database=postgres sslmode=require",
		user, password, host, port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Conn{host, port, &User{user, password}, db}, nil
}

// Close closes the Conn.
func (c *Conn) Close() error {
	return c.DB.Close()
}

// Exec execs a a series of Sequence objects against the connection.
func (c *Conn) Exec(xs ...Sequence) (err error) {
	for _, x := range xs {
		for _, s := range x.SQL() {
			_, err = c.DB.Exec(s)
			if err != nil {
				return fmt.Errorf("%v: %v", x.String(), err)
			}
		}
	}
	return
}

// CreateDatabase creates a named database and owner, writer, and reader users.
// The name will be truncated to 63 bytes, and then again to 56 bytes for the
// generated user names.
func (c *Conn) CreateDatabase(name string) (dd *DatabaseDescriptor, err error) {
	name = truncateBytes(name, 63)
	uname := truncateBytes(name, 56) // 63 - 7 to fit user names

	pw, err := genPasswords(3, 30)
	if err != nil {
		return
	}

	dd = &DatabaseDescriptor{
		Host:     c.Host,
		Port:     c.Port,
		Database: &Database{name},
		Admin:    &User{uname + "-admin", pw[0]},
		Writer:   &User{uname + "-writer", pw[1]},
		Reader:   &User{uname + "-reader", pw[2]},
	}

	err = c.execCreateDatabase(dd)
	if err != nil {
		return
	}

	c2, err := NewConn(c.Port, c.Host, dd.Admin.Name, dd.Admin.Password)
	if err != nil {
		return
	}

	err = c2.execCreateDatabaseUserPrivs(dd)
	return
}

func (c *Conn) execCreateDatabase(dd *DatabaseDescriptor) error {
	return c.Exec(
		dd.Database,
		&RevokeAllPublic{dd.Database},
		dd.Admin,
		&GrantAdmin{dd.Database, dd.Admin},
		dd.Writer,
		dd.Reader)
}

func (c *Conn) execCreateDatabaseUserPrivs(dd *DatabaseDescriptor) error {
	return c.Exec(
		&RevokeAllPublic{dd.Database},
		&GrantAccess{dd.Database, dd.Writer},
		&GrantWrite{dd.Writer},
		&GrantAccess{dd.Database, dd.Reader},
		&GrantRead{dd.Reader})
}

type Sequence interface {
	fmt.Stringer
	SQL() []string
}

// User is a Postgres user.
type User struct {
	Name     string
	Password string
}

// SQL returns the command to create this user.
func (u *User) SQL() []string {
	return []string{
		fmt.Sprintf("CREATE USER %q WITH ENCRYPTED PASSWORD '%v'", u.Name, u.Password),
	}
}

// String returns a string suitable for error messages.
func (u *User) String() string {
	return fmt.Sprintf("create user %v", u.Name)
}

// Database is a Postgres database.
type Database struct {
	Name string
}

// SQL returns the command to create this database.
func (d *Database) SQL() []string {
	return []string{
		fmt.Sprintf("CREATE DATABASE %q", d.Name),
	}
}

// String returns a string suitable for error messages.
func (d *Database) String() string {
	return fmt.Sprintf("create database %v", d.Name)
}

// GrantAccess is a grant of access on a database to a user.
type GrantAccess struct {
	On *Database
	To *User
}

// SQL returns the command to create this grant.
func (g *GrantAccess) SQL() []string {
	return []string{
		fmt.Sprintf("GRANT CONNECT ON DATABASE %q TO %q", g.On.Name, g.To.Name),
		fmt.Sprintf("GRANT USAGE ON SCHEMA PUBLIC TO %q", g.To.Name),
	}
}

// String returns a string suitable for error messages.
func (g *GrantAccess) String() string {
	return fmt.Sprintf("grant access on %v to %v", g.On.Name, g.To.Name)
}

// GrantAdmin is a grant of admin on a database to a user.
type GrantAdmin struct {
	On *Database
	To *User
}

// SQL returns the command to create this grant.
func (g *GrantAdmin) SQL() []string {
	return []string{
		fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %q TO %q WITH GRANT OPTION", g.On.Name, g.To.Name),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON SCHEMA PUBLIC TO %q WITH GRANT OPTION", g.To.Name),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA PUBLIC TO %q WITH GRANT OPTION", g.To.Name),
	}
}

// String returns a string suitable for error messages.
func (g *GrantAdmin) String() string {
	return fmt.Sprintf("grant admin on %v to %v", g.On.Name, g.To.Name)
}

// GrantRead is a grant of read to a user.
type GrantRead struct {
	To *User
}

// SQL returns the command to create this grant.
func (g *GrantRead) SQL() []string {
	return []string{
		fmt.Sprintf("GRANT SELECT ON ALL TABLES IN SCHEMA PUBLIC TO %q", g.To.Name),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA PUBLIC GRANT SELECT ON TABLES TO %q", g.To.Name),
	}
}

// String returns a string suitable for error messages.
func (g *GrantRead) String() string {
	return fmt.Sprintf("grant read to %v", g.To.Name)
}

// GrantWrite is a grant of write on a database to a user.
type GrantWrite struct {
	To *User
}

// SQL returns the command to create this grant.
func (g *GrantWrite) SQL() []string {
	priv := "SELECT,INSERT,UPDATE,DELETE,REFERENCES"
	return []string{
		fmt.Sprintf("GRANT %v ON ALL TABLES IN SCHEMA PUBLIC TO %q", priv, g.To.Name),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA PUBLIC TO %q", g.To.Name),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA PUBLIC GRANT %v ON TABLES TO %q", priv, g.To.Name),
	}
}

// String returns a string suitable for error messages.
func (g *GrantWrite) String() string {
	return fmt.Sprintf("grant write to %v", g.To.Name)
}

// RevokeAllPublic is a revocation of all privs in the public schema.
type RevokeAllPublic struct {
	On *Database
}

// SQL returns the commands to produce this revocation.
func (r *RevokeAllPublic) SQL() []string {
	return []string{
		"REVOKE ALL PRIVILEGES ON SCHEMA PUBLIC FROM PUBLIC CASCADE",
		fmt.Sprintf("REVOKE ALL PRIVILEGES ON DATABASE %q FROM PUBLIC CASCADE", r.On.Name),
	}
}

// String returns a string suitable for error messages.
func (r *RevokeAllPublic) String() string {
	return fmt.Sprintf("revoke all public on %v", r.On.Name)
}

// truncateBytes truncates string s to length n bytes. The returned string may be
// shorter than n bytes if a rune is bisected.
func truncateBytes(s string, n int) string {
	for len(s) > n {
		_, i := utf8.DecodeLastRuneInString(s)
		s = s[:len(s)-i]
	}
	return s
}

// genPassword generates a password of length l
func genPassword(l int) (p string, err error) {
	// url non-encode characters: [0-9a-zA-Z$.+!*(),_-]
	allow := &unicode.RangeTable{
		R16: []unicode.Range16{
			{0x0030, 0x0039, 1}, // '0' - '9'
			{0x0021, 0x0021, 1}, // '!' - '!'
			{0x0024, 0x0024, 1}, // '$' - '$'
			{0x0028, 0x002e, 1}, // '(' - '.'
			{0x0041, 0x005a, 1}, // 'A' - 'Z'
			{0x005f, 0x005f, 1}, // '_' - '_'
			{0x0061, 0x007a, 1}, // 'a' - 'z'
		},
	}

	strip := func(r rune) rune {
		if unicode.Is(allow, r) {
			return r
		}
		return -1
	}

	for len(p) < l {
		b := make([]byte, l*2)
		_, err = rand.Read(b)
		if err != nil {
			err = fmt.Errorf("generate password: %v", err)
			return
		}

		p += strings.Map(strip, string(b))
	}
	p = p[:l]
	return
}

// genPasswords generates n passwords of length l
func genPasswords(n, l int) (p []string, err error) {
	p = make([]string, n)
	for i := 0; i < len(p); i++ {
		p[i], err = genPassword(l)
		if err != nil {
			return
		}
	}
	return
}
