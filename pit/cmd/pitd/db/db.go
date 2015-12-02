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
	rows, err := db.Query("SELECT projects.ID, projects.Name, projects.Description, members.Email FROM projects, members WHERE projects.Owner = members.ID")
	if err != nil {
		return []Project{}, err
	}

	result := []Project{}

	for rows.Next() {
		project := Project{}

		if err = rows.Scan(&project.id, &project.name, &project.description, &project.owner); err != nil {
			return []Project{}, err
		}

		result = append(result, project)
	}

	return result, nil
}

// FindByOwner retrieves a list of all projects owned by the specified owner
func (Projects) FindByOwner(ctx context.Context, owner string) ([]Project, error) {
	var ownerid int

	err := db.QueryRow("SELECT ID from members where Email = ?;", owner).Scan(&ownerid)
	if err != nil {
		return []Project{}, err
	}

	rows, err := db.Query("SELECT * from projects where Owner = ?;", ownerid)
	if err != nil {
		return []Projects{}, err
	}

	result := []Project{}

	for rows.Next() {
		project, err := projectFromRow(rows)
		if err != nil {
			return []Project{}, err
		}

		result = append(result, project)
	}

	return result, nil
}

// FindByID retrieves the project with the given ID
func (Projects) FindByID(ctx context.Context, id int) (Project, error) {
	result := Project{}

	row := db.QueryRow("SELECT projects.ID, projects.Name, projejcts.Description, members.Email FROM projects, members WHERE projects.Owner = members.ID AND projects.ID = ?;", id)
	if err := row.Scan(&result.id, &result.name, &result.description, &result.owner); err != nil {
		return Project{}, err
	}

	return result, nil
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
