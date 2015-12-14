// Package db provides the pit database implementation
package db

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/net/context"

	_ "github.com/cznic/ql/driver"
)

const (
	Closed   = "CLOSED"
	Open     = "OPEN"
	Returned = "RETURNED"
)

var (
	ErrIssueExists           = errors.New("Issue already exists")
	ErrMemberExists          = errors.New("Member already exists")
	ErrNonContributingMember = errors.New("Member is not a project contributor")
	ErrNoSuchIssue           = errors.New("No such issue")
	ErrNoSuchMember          = errors.New("No such member")
	ErrNoSuchOwner           = errors.New("No such owner")
	ErrProjectExists         = errors.New("Project exists")
)

var (
	qldb *sql.DB
	IDs  chan int
)

func init() {
	qldb, err := sql.Open("ql", "memory://pit.db")
	if err != nil {
		panic(fmt.Sprintf("cannot create database instance: [%+v]", err))
	}

	if err = qldb.Ping(); err != nil {
		panic(fmt.Sprintf("database ping failed: [%+v]", err))
	}

	if err = mkTables(qldb); err != nil {
		panic(fmt.Sprintf("table creation failed: [%+v]", err))
	}

	go nextID(IDs)
}

func databaseFromContext(ctx context.Context) *sql.DB {
	result, ok := ctx.Value("database").(*sql.DB)
	if ok {
		return result
	}
	return qldb
}

func idsChanFromContext(ctx context.Context) chan int {
	result, ok := ctx.Value("ids-chan").(chan int)
	if ok {
		return result
	}
	return IDs
}

func mkTables(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Entity tables: projects, issues, members
	if _, err = tx.Exec("CREATE TABLE projects (ID int, Name string, Description string, Owner int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE issues (ID int, Name string,  Description string, Priority int, Status string, Project int, Reporter int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE members (ID int, Email string);"); err != nil {
		tx.Rollback()
		return err
	}

	// Association tables: contributors, assignments, watchers
	if _, err = tx.Exec("CREATE TABLE contributors (PID int, MID int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE assignments (MID int, IID int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE watchers (MID int, IID int);"); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func nextID(ch chan<- int) {
	next := 100

	for {
		ch <- next
		next++
	}
}

// test support functions

func createdb(dbname string) *sql.DB {
	db, err := sql.Open("ql", fmt.Sprintf("memory://%s.db", dbname))
	if err != nil {
		panic(fmt.Sprintf("cannot create database instance: [%+v]", err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("database ping failed: [%+v]", err))
	}

	if err = mkTables(db); err != nil {
		panic(fmt.Sprintf("table creation failed: [%+v]", err))
	}

	return db
}

func dropdb(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to drop the database: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE projects"); err != nil {
		panic(fmt.Sprintf("cannot drop the projects table: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE issues"); err != nil {
		panic(fmt.Sprintf("cannot drop the issues table: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE members"); err != nil {
		panic(fmt.Sprintf("cannot drop the members table: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE contributors"); err != nil {
		panic(fmt.Sprintf("cannot drop the contributors table: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE assignments"); err != nil {
		panic(fmt.Sprintf("cannot drop the assignments table: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE watchers"); err != nil {
		panic(fmt.Sprintf("cannont drop the watchers table: [%+v]", err))
	}

	tx.Commit()
}

func emptytables() context.Context {
	db := createdb("emptytables")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	return ctx
}

func sameissues(a, b []Issue) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	case len(a) != len(b):
		return false
	default:
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}
func samemembers(a, b []Member) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	case len(a) != len(b):
		return false
	default:
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}

func sameprojects(a, b []Project) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	case len(a) != len(b):
		return false
	default:
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}
