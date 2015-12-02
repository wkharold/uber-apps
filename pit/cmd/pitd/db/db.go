// Package db provides the pit database implementation
package db

import (
	"database/sql"
	"fmt"

	"golang.org/x/net/context"

	_ "github.com/cznic/ql"
)

var (
	db  *sql.DB
	IDs chan<- int
)

type Project struct {
	id          int
	name        string
	description string
	owner       string
}

type Projects struct{}

func init() {
	db, err := sql.Open("ql", "memory://pit.db")
	if err != nil {
		panic(fmt.Sprintf("cannot create database instance: [%+v]", err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("database ping failed: [%+v]", err))
	}

	if err = mkTables(); err != nil {
		panic(fmt.Sprintf("table creation failed: [%+v]", err))
	}

	go nextID(IDs)
}

// FindAll retrieves a list of all the projects in the repository.
func (Projects) FindAll(ctx context.Context) ([]Project, error) {
	rows, err := db.Query("SELECT * from projects;")
	if err != nil {
		return []Project{}, err
	}

	result := []Project{}

	for rows.Next() {
		project, err := projectFromRow(rows)
		if err != nil {
			return []Project{}, err
		}

		result = append(result, project)
	}
}

func projectFromRow(rows *sql.Rows) (Project, error) {
	var ownerid int
	project := Project{}

	if err := rows.Scan(&project.name, &project.description, &ownerid); err != nil {
		return Project{}, err
	}

	err := db.QueryRow("Select Email from members where ID = ?;", ownerid).Scan(&project.owner)
	if err != nil {
		return Project{}, err
	}

	return project, nil
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

	if _, err = tx.Exec("CREATE TABLE issues (ID int, Description string, Priority int, Reporter int);"); err != nil {
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
