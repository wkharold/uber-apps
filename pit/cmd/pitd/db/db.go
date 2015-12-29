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
	ErrNoSuchProject         = errors.New("No such project")
	ErrProjectExists         = errors.New("Project exists")
)

var (
	qldb *sql.DB
	IDs  chan int
)

func init() {
	var err error

	qldb, err = sql.Open("ql", "memory://pit.db")
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
